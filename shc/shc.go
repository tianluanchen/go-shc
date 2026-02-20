package shc

import (
	"bytes"
	"os"
)

var Version = "0.0.0"

type PackOption struct {
	// Specify the Customized go code generation function
	CustomGenerator func(script string) string
	GenerateOption
	CompileOption
}

// Package shell scripts to generate binary programs. If the shell option is empty, it will be automatically set
func PackShellScript(script string, opt PackOption) (*os.File, error) {
	if opt.Shell == "" {
		_, shell, index := ParseShebang(script)
		if index > 0 {
			opt.Shell = shell
		} else {
			opt.Shell = "bash"
		}
	}
	var buf bytes.Buffer
	if opt.CustomGenerator != nil {
		buf.WriteString(opt.CustomGenerator(script))
	} else {
		if err := GenerateGoCodes(script, &buf, opt.GenerateOption); err != nil {
			return nil, err
		}
	}
	return Compile(&buf, opt.CompileOption)
}
