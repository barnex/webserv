package webserv

import (
	"io"
	"net/http"
	"os"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {
	path := "." + r.URL.Path
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
