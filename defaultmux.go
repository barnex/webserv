package webserv

import (
	"fmt"
	"net/http"
)

var specialHandlers = make(map[string]http.HandlerFunc)

// main handler, dispatches to others
func defaultMux(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			Log("PANIC", err)
			//panic(err)  // for debug. TODO: log stack
		}
	}()

	logHandler(w, r) // does nothing else but logging the request

	if h, ok := specialHandlers[r.URL.Path]; ok {
		h(w, r)
	} else {
		fileHandler(w, r)
	}
}

func SetHandler(url string, f http.HandlerFunc) {
	if url[0:1] != "/" {
		fatalErr("setHandler needs absolute path: " + url)
	}
	if _, ok := specialHandlers[url]; ok {
		Log("a handler for", url, "is already present")
		return
	}
	specialHandlers[url] = f
}
