package shc

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Generate unique strings that can be used as filenames
func UniqueString() string {
	timestamp := time.Now().UnixMilli()
	chars := []byte{}
	r := rand.New(rand.NewSource(timestamp))
	for range 3 {
		chars = append(chars, byte(97+r.Intn(26)))
	}
	return fmt.Sprintf("%d%s", timestamp, string(chars))
}

// Parses into os and arch e.g. windows/amd64 => windows,amd64; linux => linux
func ParseOsarch(s string) (string, string) {
	slices := strings.Split(s, "/")
	if len(slices) >= 2 {
		return slices[0], slices[1]
	} else {
		return slices[0], ""
	}
}

// Parses shebang header, returns shebang, shell, shebang line index (counts from 0, returns -1 if none)
func ParseShebang(s string) (string, string, int) {
	slices := strings.Split(s, "\n")
	for i, v := range slices {
		v = strings.Trim(v, " \t\r")
		if strings.HasPrefix(v, "#!") {
			v := strings.TrimLeft(v, "#!")
			count := 0
			shebang, shell, line := "#!"+v, "", i
			for vv := range strings.SplitSeq(v, " ") {
				if vv == "" {
					continue
				}
				if count <= 1 {
					shell = vv
				} else {
					break
				}
				count += 1
			}
			return shebang, shell, line
		}
	}
	return "", "", -1
}
