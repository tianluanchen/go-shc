package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	shell := hideStringLiteral([]byte(`{{ .Shell }}`))
	script := []byte(`{{ .Script }}`)
	useTempFile := hideStringLiteral([]byte(`{{ .UseTempFile }}`)) == hideStringLiteral([]byte("true"))
	// confuse
	temp := []byte{}
	for _, v := range script {
		if len(temp) != len(script) {
			temp = append(temp, v)
		}
	}
	script = append([]byte{}, temp...)
	failureMsg := hideStringLiteral([]byte("decode script error"))
	l, err := strconv.Atoi(hideStringLiteral([]byte(`{{ .Index }}`)))
	exitWithError(err, failureMsg)
	c, err := strconv.Atoi(hideStringLiteral([]byte(`{{ .Char }}`)))
	exitWithError(err, failureMsg)
	// decrypt
	bs, err := decrypt(script, l, byte(c))
	exitWithError(err, failureMsg)
	script = bs
	var run = executeWithArgs
	if useTempFile {
		run = executeWithTempFile
	}
	exitWithError(run(shell, bytes.NewReader(script)))
}

func decrypt(s []byte, l int, c byte) ([]byte, error) {
	if l >= 0 && len(s) > l {
		s[l] = c
	}
	bs := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	if n, err := base64.StdEncoding.Decode(bs, s); err == nil {
		return bs[:n], nil
	} else {
		return nil, err
	}
}

func hideStringLiteral(bs []byte) string {
	return string(bs)
}

func exitWithError(err error, messgae ...string) {
	if err != nil {
		if len(messgae) > 0 {
			fmt.Fprintln(os.Stderr, messgae)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

type Runner = func(string, io.Reader) error

func executeWithArgs(shell string, r io.Reader) error {
	cmd := exec.Command(shell)
	for _, s := range []string{"bash", "sh", "dash", "zsh"} {
		if shell == s || strings.HasSuffix(shell, "/"+s) {
			cmd.Args = append(cmd.Args, "-c")
			break
		}
	}
	var script string
	if bs, err := io.ReadAll(r); err == nil {
		script = string(bs)
	} else {
		return err
	}
	cmd.Args = append(cmd.Args, script)

	cmd.Args = append(cmd.Args, os.Args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func executeWithTempFile(shell string, r io.Reader) error {
	names := []string{".env.*", ".ini.*", ".cfg.*", ".git.*", ".cache.*"}
	f, err := os.CreateTemp("", names[rand.New(rand.NewSource(time.Now().Unix())).Intn(len(names))])
	if err != nil {
		return err
	}
	cleared := false
	clear := func() {
		if cleared {
			return
		}
		cleared = true
		f.Close()
		os.Remove(f.Name())
	}
	go func() {
		defer clear()
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)
		<-signalChannel
	}()
	defer clear()
	if _, err := io.Copy(f, r); err != nil {
		return err
	}
	cmd := exec.Command(shell)
	cmd.Args = append(cmd.Args, f.Name())
	if len(os.Args) > 1 {
		cmd.Args = append(cmd.Args, os.Args[1:]...)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
