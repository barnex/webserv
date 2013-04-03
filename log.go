package webserv

import (
	"log"
	"net/http"
)

var requestcount int64

// log the request
func logHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	requestcount++
}
