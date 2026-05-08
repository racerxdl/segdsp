package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed content/index.html
var indexHTML []byte

//go:embed content/assets
var assetsFS embed.FS

func content(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(indexHTML)
		return
	}
	assetsServer.ServeHTTP(w, r)
}

var assetsServer http.Handler

func init() {
	sub, _ := fs.Sub(assetsFS, "content/assets")
	assetsServer = http.StripPrefix("/assets/", http.FileServer(http.FS(sub)))
}
