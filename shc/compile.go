package shc

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

type CompileOption struct {
	TempDir    string // Directory for temporary storage of source code
	Output     string // Output path, if empty, use OutputDir
	OutputDir  string // Output directory, if empty, use TempDir
	Osarch     string // Target Os and Arch
	TrimPath   bool   // Whether or not to exclude the source path when compiling.
	GoCompiler string // Compiler name or path, default "go"
}

// Compile a single go code file and disable the module
func Compile(r io.Reader, opt CompileOption) (*os.File, error) {
	goos, goarch := ParseOsarch(opt.Osarch)
	if opt.GoCompiler == "" {
		opt.GoCompiler = "go"
	}
	if opt.TempDir == "" {
		opt.TempDir = os.TempDir()
	}
	if opt.Output == "" {
		if opt.OutputDir == "" {
			opt.OutputDir = opt.TempDir
		}
		opt.Output = path.Join(opt.OutputDir, UniqueString())
		if runtime.GOOS == "windows" {
			opt.Output += ".exe"
		}
	}
	f, err := os.CreateTemp(opt.TempDir, UniqueString()+"*.go")
	if err != nil {
		return nil, err
	}
	cleared := false
	clear := func() {
		if cleared {
			return
		}
		f.Close()
		os.Remove(f.Name())
		cleared = true
	}
	defer clear()
	if _, err := io.Copy(f, r); err != nil {
		return nil, err
	}
	cmd := exec.Command(opt.GoCompiler, "build", "-ldflags", "-s -w")
	if opt.TrimPath {
		pwd, err := filepath.Abs(opt.TempDir)
		if err != nil {
			return nil, err
		}
		cmd.Args = append(cmd.Args, "-gcflags", "all=-trimpath="+pwd, "-asmflags", "all=-trimpath="+pwd)
	}
	cmd.Args = append(cmd.Args, "-o", opt.Output, f.Name())
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, "GO111MODULE=off")
	if goos != "" {
		cmd.Env = append(cmd.Env, "GOOS="+goos)
	}
	if goarch != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+goarch)
	}
	if output, err := cmd.CombinedOutput(); err != nil {
		if len(output) > 0 {
			return nil, errors.New(string(output))
		} else {
			return nil, err
		}
	}
	clear()
	return os.Open(opt.Output)
}
