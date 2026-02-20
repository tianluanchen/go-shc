package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"go-shc/shc"
	"go-shc/web"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a http server to provide services",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		trimPath, _ := cmd.Flags().GetBool("trim-path")
		fmt.Println("Server listening on", addr)
		exitWithError(http.ListenAndServe(addr, NewApp(trimPath)))
	},
}

func init() {
	serveCmd.Flags().Bool("trim-path", false, "Trim source code path")
	serveCmd.Flags().StringP("addr", "a", "127.0.0.1:8080", "Listen address")
	rootCmd.AddCommand(serveCmd)
}

func NewApp(trimPath bool) *http.ServeMux {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(web.FS))
	etag := fmt.Sprintf(`"%s"`, hex.EncodeToString(web.MD5Hash))
	respond := func(w http.ResponseWriter, statusCode int, msg string) {
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
	}
	notfound := func(w http.ResponseWriter) {
		respond(w, http.StatusNotFound, "404 Not Found")
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, ok := web.PathMap[r.URL.Path]; ok {
			w.Header().Set("ETag", etag)
			fileServer.ServeHTTP(w, r)
			return
		}
		notfound(w)
	})
	mux.HandleFunc("/shc", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			notfound(w)
			return
		}
		if r.ContentLength >= 5<<20 {
			respond(w, http.StatusBadRequest, "The request body is too large")
			return
		}
		q := r.URL.Query()
		if q.Get("osarch") == "" {
			respond(w, http.StatusBadRequest, "The osarch is required")
			return
		}
		opt := shc.PackOption{
			GenerateOption: shc.GenerateOption{
				Shell:       q.Get("shell"),
				UseTempFile: len(q.Get("useTempFile")) > 0,
			},
			CompileOption: shc.CompileOption{
				Osarch:   q.Get("osarch"),
				TrimPath: trimPath,
			},
		}
		var body io.Reader
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			r.ParseMultipartForm(1 << 20)
			if s := r.FormValue("script"); s == "" {
				if f, _, err := r.FormFile("script"); err != nil {
					respond(w, http.StatusBadRequest, "The script file is required")
					return
				} else {
					defer f.Close()
					body = f
				}

			} else {
				body = strings.NewReader(s)
			}

		} else {
			body = r.Body
		}
		bs, err := io.ReadAll(body)
		if err != nil {
			respond(w, http.StatusInternalServerError, fmt.Sprintf("Read body error: %v", err))
			return
		}
		f, err := shc.PackShellScript(string(bs), opt)
		if err != nil {
			respond(w, http.StatusInternalServerError, fmt.Sprintf("Failed to compile: %v", err))
			return
		} else {
			defer os.Remove(f.Name())
			defer f.Close()
			if strings.Contains(opt.Osarch, "windows") {
				w.Header().Set("Content-Disposition", "attachment; filename=app.exe")
			} else {
				w.Header().Set("Content-Disposition", "attachment; filename=app")
			}
			io.Copy(w, f)
		}
	})
	return mux
}
