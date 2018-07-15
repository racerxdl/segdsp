package main

import (
	"net/http"
)

func content(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "content/index.html")
}
