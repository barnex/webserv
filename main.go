package webserv

import (
	"flag"
	"net/http"
)

var flag_http = flag.String("http", ":80", "http port")

// After setup, run the web server on the port specified by -http flag.
func Run() {
	flag.Parse()
	http.HandleFunc("/", defaultMux)
	Log("Listen and serve", *flag_http)
	fatalErr(http.ListenAndServe(*flag_http, nil))
}
