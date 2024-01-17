package web

import (
	"embed"
	"io"
	"io/fs"
	"strings"
)

//go:embed static/*
var staticFS embed.FS

// static resource
var FS fs.FS

// Map of the url path
var PathMap = map[string]string{}

func init() {
	var err error
	FS, err = fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	// url root path
	if ExistFile("index.html") {
		PathMap["/"] = "index.html"
	}
	Walk(".", func(path string, info fs.FileInfo, data []byte) error {
		urlPath := "/" + path
		if info.IsDir() {
			urlPath := urlPath + "/"
			indexPage := path + "/index.html"
			if ExistFile(indexPage) {
				PathMap[urlPath] = indexPage
			}
		} else {
			PathMap[urlPath] = path
		}
		return nil
	})
}

// Determine if a file exists in a path
func ExistFile(path string) bool {
	path = strings.TrimLeft(path, "/")
	f, err := FS.Open(path)
	if err != nil {
		return false
	}
	info, _ := f.Stat()
	return !info.IsDir()
}

// Iterate through all the files in the static resource
func Walk(root string, fn func(path string, info fs.FileInfo, data []byte) error) error {
	return fs.WalkDir(FS, root, func(path string, d fs.DirEntry, err error) error {
		f, _ := FS.Open(path)
		info, _ := f.Stat()
		var bs []byte
		if !info.IsDir() {
			bs, _ = io.ReadAll(f)
		}
		return fn(path, info, bs)
	})
}
