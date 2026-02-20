/**
 * @typedef {Object} ObfuscateOption
 * @property {string} shell - Specify the shell.
 * @property {string} varNamePrefix - Variable name prefix.
 * @property {number} sliceLength - Slice length.
 * @property {number} minVarNameLength - Minimum variable name length.
 * @property {boolean} useTempFile - Create a temporary file for execution.
 * @property {number} argLengthLimit - Argument length limit; exceeding this limit will force script execution through the creation of a temporary file.
 * @property {number} varCountLimit - Variable count limit; exceeding this limit will automatically set the slice length.
 */

/** 
 * @param {string} codes
 * @param {ObfuscateOption} option default {
        sliceLength: 5,
        minVarNameLength: 3,
        varNamePrefix: "_",
        useTempFile: false,
        argLengthLimit: 17200,
        varCountLimit: 20000,
    }
*/
async function obfuscateShellScript(codes, option = {}) {
    const defaultOption = {
        sliceLength: 5,
        minVarNameLength: 3,
        varNamePrefix: "_",
        useTempFile: false,
        argLengthLimit: 17200,
        varCountLimit: 20000,
    };
    let { shell, sliceLength, varNamePrefix, minVarNameLength, useTempFile, argLengthLimit, varCountLimit } = Object.assign({}, defaultOption, option);
    const lines = codes.split("\n");
    const shebang = [-1, ""];
    for (let i = 0; i < lines.length; i++) {
        const l = lines[i].trim();
        if (l !== "") {
            if (l.startsWith("#!")) {
                shebang[0] = i;
                shebang[1] = l;
            }
            break;
        }
    }
    if (shebang[0] > -1) {
        lines.splice(shebang[0], 1);
        codes = lines.join("\n");
    }
    if (["", null, undefined].includes(shell)) {
        if (shebang[0] > -1) {
            shell = "";
            const ss = shebang[1]
                .split(" ")
                .map(e => e.trim())
                .filter(e => e !== "");
            if (ss.length > 0) {
                if (ss.length === 1) {
                    shell = ss[0].replace("#!", "");
                } else {
                    shell = ss[1];
                }
            }
        } else {
            shell = "bash";
        }
    }
    if (shell === "") {
        shell = "bash";
    }
    /**
     * @param {string} s
     * @return {Promise<string>}
     */
    const b64encode = s => {
        return new Promise((resolve, reject) => {
            if (typeof Buffer === "function") {
                resolve(Buffer.from(s).toString("base64"));
                return;
            }
            const blob = new Blob([s], { type: "text/plain" });
            const reader = new FileReader();
            reader.onload = () => {
                resolve(reader.result.replace("data:text/plain;base64,", ""));
            };
            reader.onerror = reject;
            reader.readAsDataURL(blob);
        });
    };
    const chars = Array(26)
        .fill(null)
        .map((_, i) => {
            return String.fromCharCode(65 + i);
        })
        .concat(
            Array(26)
                .fill(null)
                .map((_, i) => {
                    return String.fromCharCode(97 + i);
                })
        );
    const genVarName = (length = 1) => {
        let v = "";
        for (let i = 0; i < length; i++) {
            v += chars[Math.floor(Math.random() * chars.length)];
        }
        return v;
    };
    const b64Codes = await b64encode(codes);
    const tempFile = varNamePrefix + genVarName(3);
    let evalCodes = "";
    // If the maximum limit on the length of command line arguments is exceeded, a temporary file is used
    useTempFile = useTempFile || b64Codes.length > argLengthLimit;
    if (useTempFile) {
        evalCodes = `trap "rm -f \\$${tempFile}" EXIT;${tempFile}=$(mktemp) || exit 1; echo ${b64Codes} | base64 -d > "$${tempFile}";${shell} "$${tempFile}" "$@";`;
    } else {
        evalCodes = `${shell} -c "$(base64 -d <<< "${b64Codes}")" ${shell} "$@"`;
    }
    // variable count
    let varCount = Math.ceil(evalCodes.length / sliceLength);
    // setting a maximum limit on the number of variables to prevent execution overflow
    if (varCount > varCountLimit) {
        varCount = varCountLimit;
        sliceLength = Math.ceil(evalCodes.length / varCount);
    }
    // variable name length
    let varNameLength = minVarNameLength;
    while (chars.length ** varNameLength < varCount) {
        varNameLength++;
    }
    // generate
    let tmp = evalCodes;
    let varNameSet = new Set();
    let varGroups = [];
    while (tmp.length > 0) {
        const value = tmp.slice(0, sliceLength);
        tmp = tmp.slice(sliceLength);
        let key = "";
        do {
            key = genVarName(varNameLength);
        } while (varNameSet.has(key));
        varNameSet.add(key);
        // prefix
        key = varNamePrefix + key;
        varGroups.push({ key, value });
    }

    let result = "";
    if (shebang[0] >= 0) {
        result += shebang[1] + "\n";
    }
    // Disrupt order
    for (const { key, value } of varGroups.map(e => e).sort(() => Math.random() - 0.5)) {
        result += `${key}='${value}';`;
    }
    result += `eval "${varGroups.map(e => "$" + e.key).join("")}";`;
    return {
        source: codes,
        result: result,
        useTempFile,
        varCount,
        sliceLength,
    };
}
