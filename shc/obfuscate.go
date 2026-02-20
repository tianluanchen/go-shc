package shc

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type ObfuscateOption struct {
	// Specify a shell that is automatically parsed and configured by default
	Shell string
	// Variable prefix, default "_"
	VarNamePrefix string
	// Slice length, default 5
	SliceLength int
	// Minimum variable length, default 3
	MinVarNameLength int
	// Run by creating a temporary file
	UseTempFile bool
	// Argument length limit; exceeding this limit will force script execution through the creation of a temporary file. default 17200
	ArgLengthLimit int
	// Variable count limit; exceeding this limit will automatically set the slice length. default 20000
	VarCountLimit int
}

type ObfuscateResult struct {
	Output string
	ObfuscateOption
}

func ObfuscateShellScript(script string, opt ObfuscateOption) *ObfuscateResult {
	if opt.SliceLength < 1 {
		opt.SliceLength = 5
	}
	if opt.MinVarNameLength < 1 {
		opt.MinVarNameLength = 3
	}
	if opt.ArgLengthLimit < 1 {
		opt.ArgLengthLimit = 17200
	}
	if opt.VarCountLimit < 1 {
		opt.VarCountLimit = 20000
	}
	if opt.VarNamePrefix == "" {
		opt.VarNamePrefix = "_"
	}
	shebang, shell, line := ParseShebang(script)
	if opt.Shell == "" {
		if line < 0 {
			opt.Shell = "bash"
		} else {
			opt.Shell = shell
		}
	}
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	genVarName := func(length int) string {
		bs := []byte{}
		for i := 0; i < length; i++ {
			bs = append(bs, chars[r.Intn(len(chars))])
		}
		return string(bs)
	}
	b64Codes := base64.StdEncoding.EncodeToString([]byte(script))
	tempFile := opt.VarNamePrefix + genVarName(3)
	evalCodes := ""
	// If the maximum limit on the length of command line arguments is exceeded, a temporary file is used
	opt.UseTempFile = opt.UseTempFile || len(b64Codes) > opt.ArgLengthLimit
	if opt.UseTempFile {
		evalCodes = fmt.Sprintf(`trap "rm -f \$%s" EXIT;%s=$(mktemp) || exit 1; echo %s | base64 -d > "$%s";%s "$%s" "$@";`, tempFile, tempFile, b64Codes, tempFile, opt.Shell, tempFile)
	} else {
		evalCodes = fmt.Sprintf(`%s -c "$(base64 -d <<< "%s")" %s "$@"`, opt.Shell, b64Codes, opt.Shell)
	}
	// setting a maximum limit on the number of variables to prevent execution overflow
	varCount := int(math.Ceil(float64(len(evalCodes)) / float64(opt.SliceLength)))
	if varCount > opt.VarCountLimit {
		varCount = opt.VarCountLimit
		opt.SliceLength = int(math.Ceil(float64(len(evalCodes)) / float64(varCount)))
	}
	varNameLength := opt.MinVarNameLength
	for int(math.Pow(float64(len(chars)), float64(varNameLength))) < opt.VarCountLimit {
		varNameLength++
	}
	// generate
	tmp := evalCodes
	varNameMap := make(map[string]struct{})
	varGroups := [][]string{}
	for len(tmp) > 0 {
		sliceLength := min(opt.SliceLength, len(tmp))
		v := tmp[0:sliceLength]
		tmp = tmp[sliceLength:]
		k := ""
		for {
			k = genVarName(varNameLength)
			if _, ok := varNameMap[k]; !ok {
				break
			}
		}
		varNameMap[k] = struct{}{}
		k = opt.VarNamePrefix + k
		varGroups = append(varGroups, []string{k, v})
	}
	result := ""
	if line > -1 {
		result = shebang + "\n"
	}
	varNames := make([]string, 0, len(varGroups))
	for _, v := range varGroups {
		varNames = append(varNames, v[0])
	}
	// Disrupt order
	sort.Slice(varGroups, func(i, j int) bool {
		return rand.Intn(2) == 0
	})
	for _, v := range varGroups {
		result += fmt.Sprintf(`%s='%s';`, v[0], v[1])
	}
	result += `eval "`
	for _, v := range varNames {
		result += fmt.Sprintf(`$%s`, v)
	}
	result += `";`
	return &ObfuscateResult{
		Output:          result,
		ObfuscateOption: opt,
	}
}
