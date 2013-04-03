package webserv

// I/O utilities

import (
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

func loadTemplate(templ string) *template.Template {
	templText, err := ioutil.ReadFile(templ)
	fatalErr(err)
	return template.Must(template.New(templ).Parse(string(templText)))
}

// functionality of readlink unix program
func readlink(f string) string {
	l, err := os.Readlink(f)
	if err != nil {
		//Log(err)
		return f
	}
	if !path.IsAbs(l) {
		l = path.Dir(f) + "/" + l
	}
	Log("READLINK", f, "->", l)
	return l
}

// wrapper for os.Stat, returns nil on error
func stat(f string) os.FileInfo {
	i, err := os.Stat(f)
	if err != nil {
		Log(err)
	}
	return i
}
