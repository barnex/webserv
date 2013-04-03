package webserv

// File: logging and error reporting utility functions
// Author: Arne Vansteenkiste

import (
	"log"
	"path"
)

// If err != nil, trigger log.Fatal(msg, err)
func fatalErr(err interface{}, msg ...interface{}) {
	if err != nil {
		log.Fatal(append(msg, err)...)
	}
}

// Panics if err is not nil. Signals a bug.
func panicErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// Logs the error of non-nil, plus message
func logErr(err error, msg ...interface{}) {
	if err != nil {
		log.Println(append(msg, err)...)
	}
}

// Panics with "illegal argument" if test is false.
func argument(test bool) {
	if !test {
		log.Panic("illegal argument")
	}
}

// Panics with "assertion failed" if test is false.
func assert(test bool) {
	if !test {
		log.Panic("assertion failed")
	}
}

// Remove extension from file name.
func noExt(file string) string {
	ext := path.Ext(file)
	return file[:len(file)-len(ext)]
}

func Log(msg ...interface{}) {
	log.Println(msg...)
}
