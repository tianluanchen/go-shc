package shc

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"text/template"
	"time"
)

//go:embed template/main.go
var templateBody []byte
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("shc").Parse(string(templateBody)))
}

type GenerateOption struct {
	Shell       string
	UseTempFile bool
}
type tmplData struct {
	GenerateOption
	Script string
	Index  int
	Char   byte
}

// Generate Go code that wraps shell scripts
func GenerateGoCodes(script string, w io.Writer, opt GenerateOption) error {
	if opt.Shell == "" {
		return fmt.Errorf("shell is empty")
	}
	var data = tmplData{
		GenerateOption: opt,
		Script:         script,
	}
	s, i, c := encrypt([]byte(data.Script))
	data.Script = string(s)
	data.Index = i
	data.Char = c
	return tmpl.Execute(w, data)
}

// After base64 encoding, replace a random byte at a random location, return the byte slice and the index and value of the replaced byte
func encrypt(s []byte) ([]byte, int, byte) {
	bs := make([]byte, base64.StdEncoding.EncodedLen(len(s)))
	base64.StdEncoding.Encode(bs, s)
	// get length of valid part (exclude =)
	validLength := 0
	for i := len(bs) - 1; i >= 0; i-- {
		if bs[i] != 61 {
			validLength = i + 1
			break
		}
	}
	// if validLength <= 0, do not modify
	if validLength <= 0 {
		return bs, -1, 0
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	index := r.Intn(validLength)
	var c byte
	if r.Intn(2) > 0 {
		c = byte(97 + r.Int63n(26))
	} else {
		c = byte(65 + r.Int63n(26))
	}
	original := bs[index]
	bs[index] = c
	return bs, index, original
}

// Removes the bytes at the specified index before base64 decoding , this function is used in the template
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
