package shc

import (
	"os"
	"strings"
	"testing"
)

func TestEncrypt(t *testing.T) {
	bs, err := decrypt(encrypt([]byte("")))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs), len(bs))
}

func TestCompile(t *testing.T) {
	codes := `
	package main
	func main(){
	}
	`
	f, err := Compile(strings.NewReader(codes), CompileOption{})
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	t.Log(f.Name())
}

func TestPack(t *testing.T) {
	codes := `
	echo $0 $1;
	read a;
	echo $a;
	`
	opt := PackOption{}
	opt.OutputDir = "../tmp"
	f, err := PackShellScript(codes, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()
	t.Log(f.Name())
}

func TestParseShebang(t *testing.T) {
	codes := `
	#!/bin/bash 
	#!/bin/bash bash
	`
	shebang, shell, line := ParseShebang(codes)
	t.Logf("<%s>  <%s>  <%d>", shebang, shell, line)
}

func TestObfuscate(t *testing.T) {
	codes := `
	#!/usr/bin/env sh
	echo $0 $1;
	read a;echo $a
	`
	// codes = "echo 123"
	r := ObfuscateShellScript(codes, ObfuscateOption{
		// SliceLength: 1000,
		// UseTempFile: true,
	})
	t.Log(r.Output)
}
