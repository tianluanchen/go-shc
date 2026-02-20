package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// If err is not nil, exit with error
func exitWithError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getGlobFiles(patterns ...string) []string {
	files := []string{}
	for _, p := range patterns {
		matches, err := filepath.Glob(p)
		if err == nil {
			files = append(files, matches...)
		}
	}
	return files
}
