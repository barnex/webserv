package webserv

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
	"time"
)

var rootTemplate *template.Template

var rootNode *Content

// The root template (homedir/template.html, with homedir the directory passed to LoadContent)
// is executed with this data structure as input. Thus, the fields and methods of this struct
// can be used in the root template. E.g.:
// 	{{.Title}} {{.Uptime}}
type Content struct {
	Title    string
	Href     string
	Body     string
	Children []*Content
}

func newContent(title, href, body string) *Content {
	return &Content{title, href, body, []*Content{}}
}

// Loads and makes available all content in homedir and subdirectories.
func LoadContent(homedir string) {
	templ := "template.html"
	rootTemplate = loadTemplate(templ)
	rootNode = loadDir(homedir)
	Log(len(specialHandlers), "url handlers")
}

// return file contents, empty string in case of error.
func loadFile(fname string) string {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		logErr(err)
		return ""
	}
	return string(bytes)
}

// turn file name into suited page title
func titleOf(fname string) string {
	title := path.Base(fname)
	if title == "/" {
		title = "Home"
	}
	return title
}

func loadDir(p string) *Content {
	// clean path
	url := path.Clean(p)
	if url == "." {
		url = ""
	}
	url = "/" + url
	Log("preload dir", p)

	// load body
	body := loadFile(p + "/index.html")
	Content := newContent(titleOf(url), url, body)
	SetHandler(url, func(w http.ResponseWriter, r *http.Request) {
		renderContent(Content, w, r)
	})

	ls := readDir(p)
	for _, f := range ls {
		fullname := readlink(p + "/" + f.Name())
		f = stat(fullname)
		if f.IsDir() {
			Content.addChild(loadDir(fullname))
		} else {
			if strings.HasSuffix(fullname, ".html") && f.Name() != "index.html" && f.Name() != "template.html" {
				body := loadFile(fullname)
				Content.addChild(newContent("", "", body))
			}
		}
	}
	return Content
}

func renderContent(Content *Content, w http.ResponseWriter, r *http.Request) {
	fatalErr(rootTemplate.Execute(w, Content)) // TODO: not fatal, just panic and catch
}

func readDir(p string) []os.FileInfo {
	// load subdirs
	p = readlink(p)
	f, err := os.Open(p)
	if err != nil {
		Log(err)
		return nil
	}
	ls, err2 := f.Readdir(-1)
	if err2 != nil {
		Log(err2)
		return nil
	}
	return ls
}

func (c *Content) addChild(child *Content) {
	c.Children = append(c.Children, child)
}

// Returns the Go runtime version.
func (c *Content) GoVersion() string { return runtime.Version() }

var booted = time.Now()

// Returns the webserver's uptime.
func (c *Content) Uptime() string { return time.Since(booted).String() }

// Returns the number of requests served since boot.
func (c *Content) RequestCount() string { return fmt.Sprint(requestcount) }
