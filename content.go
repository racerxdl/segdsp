package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed content/index.html
var indexHTML []byte

//go:embed content/static
var staticFS embed.FS

func content(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func init() {
	sub, _ := fs.Sub(staticFS, "content/static")
	staticFileServer = http.FileServer(http.FS(sub))
}

var staticFileServer http.Handler
