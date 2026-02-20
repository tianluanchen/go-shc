const languageOptions = [
    ["English", "index.html"],
    ["中文", "zh.html"],
];
getFormElemHandler("#language")
    .setOptions(
        languageOptions,
        (languageOptions.find(e => {
            const v = e[1];
            const slice = location.pathname.split("/");
            let base = slice[slice.length - 1];
            if (base === "") {
                base = "index.html";
            }
            return base === v;
        }) || ["", "index.html"])[1]
    )
    .el.addEventListener("change", evt => {
        const v = evt.target.value;
        location.href = v === "index.html" ? "./" : v;
    });

// obfuscate
const obfuscateInput = getFormElemHandler("#obfuscateInput").persist("obfuscateInput");
const obfuscateOutput = getFormElemHandler("#obfuscateOutput");
const obfuscateUseTempFile = getFormElemHandler("#obfuscateUseTempFile").persist("obfuscateUseTempFile", false);

const obfuscateShellSelect = getFormElemHandler("#obfuscateShell").setOptions().persist("obfuscateShell");

const obfuscateLevelSelect = getFormElemHandler("#obfuscateLevel").setOptions().persist("obfuscateLevel");

const obfuscateShellCustom = getFormElemHandler("#obfuscateShellCustom").persist("obfuscateShellCustom");
obfuscateShellSelect.el.addEventListener("change", () => {
    obfuscateShellCustom.el.style.display = obfuscateShellSelect.getValue() === "custom" ? "" : "none";
});
obfuscateShellSelect.el.dispatchEvent(new Event("change"));

const fillObfuscateInputWithFile = importLocalFileWrapper(file => {
    file.text().then(s => obfuscateInput.setValue(s));
});

const generateObfuscateCodes = async () => {
    let shell = obfuscateShellSelect.getValue();
    if (shell == "custom") {
        shell = obfuscateShellCustom.getValue().trim();
        if (shell === "") {
            alert("shell can't be empty!");
            return;
        }
    }
    const [minVarNameLength, sliceLength] = obfuscateLevelSelect
        .getValue()
        .split(",")
        .map(e => parseInt(e));
    obfuscateOutput.setValue((await obfuscateShellScript(obfuscateInput.getValue(), { shell, sliceLength, minVarNameLength, useTempFile: obfuscateUseTempFile.getValue() })).result);
};

// pack
const packInput = getFormElemHandler("#packInput").persist("packInput");

const packShellSelect = getFormElemHandler("#packShell").setOptions().persist("packShell");

const packShellCustom = getFormElemHandler("#packShellCustom").persist("packShellCustom");

packShellSelect.el.addEventListener("change", () => {
    packShellCustom.el.style.display = packShellSelect.getValue() === "custom" ? "" : "none";
});
packShellSelect.el.dispatchEvent(new Event("change"));

const osarchSelect = getFormElemHandler("#osarch").setOptions().persist("osarch");
const packUseTempFile = getFormElemHandler("#packUseTempFile").persist("packUseTempFile", false);

const fillPackInputWithFile = importLocalFileWrapper(file => {
    file.text().then(s => packInput.setValue(s));
});

const getUrl = () => {
    const url = new URL("/shc", location);
    if (packUseTempFile.getValue()) {
        url.searchParams.set("useTempFile", "true");
    }
    url.searchParams.set("shell", packShellSelect.getValue() === "custom" ? packShellCustom.getValue().trim() : packShellSelect.getValue());
    url.searchParams.set("osarch", osarchSelect.getValue());
    return url.href;
};

const generateCurlCmd = () => {
    const root = document.querySelector("#curl");
    root.style.display = "";
    const spanList = [...root.querySelectorAll("span")];
    [`curl -X "POST" --data-raw '<script>' -o app "${getUrl()}"`, `curl -X "POST" --data-binary "@<script-file>" -o app "${getUrl()}"`].forEach((v, i) => {
        spanList[i].textContent = v;
        const btn = spanList[i].nextElementSibling;
        btn.onclick = () => copyWithBtn(v, btn);
    });
};

const getPackFile = async () => {
    const form = document.querySelector("#send");
    if (packInput.getValue().trim() === "") {
        alert("script can't be empty!");
        return;
    }
    const shell = packShellSelect.getValue();
    if (shell == "custom") {
        shell = packShellCustom.getValue().trim();
        if (shell === "") {
            alert("shell can't be empty!");
            return;
        }
    }
    form.action = getUrl();
    form.querySelector("textarea").value = packInput.getValue();
    form.submit();
};

document.querySelectorAll("textarea").forEach(e => {
    e.addEventListener("input", () => {
        if (!e.nextElementSibling) {
            return;
        }
        const value = e.value;
        const lines = (value.match(/\n/g)?.length || 0) + 1;
        const count = value.trim().length;
        const size = new Blob([value]).size;
        const elem = e.nextElementSibling;
        elem.textContent = (elem.dataset.tmpl || "")
            .replace("#char#", count + "")
            .replace("#line#", lines + "")
            .replace("#size#", formatSize(size));
    });
});

/**
 * @param {number} n
 */
function formatSize(n) {
    if (n < 1024) {
        return `${n}B`;
    } else if (n < 1024 * 1024) {
        return `${(n / 1024).toFixed(1)}KB`;
    } else {
        return `${(n / 1024 / 1024).toFixed(1)}MB`;
    }
}

/**
 *
 * @param {()=>Promise} p
 * @param {HTMLElement} el
 */
function asyncWrapperWithElem(p, el) {
    if (el.dataset.pending) {
        return;
    }
    el.dataset.pending = "true";
    el.disabled = true;
    Promise.resolve(p()).finally(() => {
        el.disabled = false;
        delete el.dataset.pending;
    });
}

/**
 * get handler for textarea, input:text, input:checkbox, select
 * @param {string|HTMLElement} sel
 * @returns
 */
function getFormElemHandler(sel) {
    return {
        el: sel instanceof HTMLElement ? sel : document.querySelector(sel),
        /**
         * arr item is ["label", "value"] or "value"
         */
        setOptions(arr, defaultIndex) {
            try {
                if (arr === undefined) {
                    arr = JSON.parse(this.el.dataset.options);
                }
                if (defaultIndex === undefined) {
                    defaultIndex = JSON.parse(this.el.dataset.default);
                }
            } catch (error) {
                console.error(this.el, "setOptions error:", error);
                return this;
            }
            arr.forEach(e => {
                const opt = document.createElement("option");
                if (e instanceof Array) {
                    opt.textContent = e[0];
                    opt.value = e[1];
                } else {
                    opt.value = opt.textContent = e;
                }
                this.el.appendChild(opt);
            });
            if (typeof defaultIndex !== "number") {
                defaultIndex = arr.findIndex(v => {
                    return v instanceof Array ? v[1] === defaultIndex : v === defaultIndex;
                });
                if (defaultIndex < 0) {
                    defaultIndex = 0;
                }
            }
            this.el.value = arr[defaultIndex] instanceof Array ? arr[defaultIndex][1] : arr[defaultIndex];

            return this;
        },
        persist(key, initV) {
            try {
                let v = localStorage.getItem(key);
                if (v === null) {
                    if (initV !== undefined) {
                        v = initV;
                    } else {
                        throw "";
                    }
                }
                this.setValue(JSON.parse(v));
            } catch { }
            const dump = () => localStorage.setItem(key, JSON.stringify(this.getValue()));
            if (this.el instanceof HTMLInputElement && ["password", "text"].includes(this.el.type)) {
                this.el.addEventListener("input", dump);
                this.el.dispatchEvent(new Event("input"));
            } else {
                this.el.addEventListener("change", dump);
                this.el.dispatchEvent(new Event("change"));
            }
            return this;
        },
        getValue() {
            if (this.el instanceof HTMLInputElement && this.el.type == "checkbox") {
                return this.el.checked;
            }
            return this.el.value;
        },
        setValue(v) {
            // ignore invalid value
            if (this.el instanceof HTMLSelectElement && [...this.el.querySelectorAll("option")].every(e => e.value !== v)) {
                return this;
            }

            if (this.el instanceof HTMLInputElement && this.el.type == "checkbox") {
                this.el.checked = Boolean(v);
            } else {
                this.el.value = v;
            }
            this.el.dispatchEvent(new Event("change"));
            this.el.dispatchEvent(new Event("input"));
            return this;
        },
    };
}

/**
 *
 * @param {(f:File)=>void} callback
 * @param {string} accept
 * @returns
 */
function importLocalFileWrapper(callback, accept = "*") {
    const input = document.createElement("input");
    input.type = "file";
    input.accept = accept;
    input.style.display = "none";
    document.body.appendChild(input);
    input.onchange = () => {
        const file = input.files.item(0);
        file && callback(file);
        input.value = "";
    };
    return () => input.click();
}

/**
 *
 * @param {string} s
 * @param {string} name
 */
function downloadTextFile(s, name = "text.txt") {
    const blob = new Blob([s], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = name;
    a.click();
    URL.revokeObjectURL(url);
}

/**
 * @param {string} s
 * @param {HTMLButtonElement} btn
 */
function copyWithBtn(s, btn) {
    clearTimeout(btn.dataset.timeout);
    if (btn.dataset.text === undefined) {
        btn.dataset.text = btn.textContent;
    }
    asyncWrapperWithElem(
        () =>
            copyText(s).finally(() => {
                btn.textContent = btn.dataset.tip;
                btn.dataset.timeout = setTimeout(() => {
                    btn.textContent = btn.dataset.text;
                }, 1000);
            }),
        btn
    );
}

/**
 * @param {string} e
 * @param {string|HTMLElement} t
 * @returns
 */
function copyText(e, t = "body") {
    const o = document.createElement("textarea");
    o.setAttribute("style", "opacity:0;position:fixed;top:-200px;height:5px");
    const n = (e, n) => {
        let c;
        (c = t instanceof HTMLElement ? t : document.querySelector(t)), (o.value = e), c.appendChild(o), o.select(), document.execCommand("copy"), o.remove(), n();
    };
    return new Promise(t => {
        if (window.navigator.clipboard)
            try {
                window.navigator.clipboard
                    .writeText(e)
                    .then(t)
                    .catch(o => {
                        n(e, t);
                    });
            } catch (o) {
                n(e, t);
            }
        else n(e, t);
    });
}
