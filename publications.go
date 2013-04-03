package webserv

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path"
	"regexp"
	"text/template"
)

var (
	pubs        []*pub
	pubcontent  *Content
	pubTemplate *template.Template
)

type pub struct {
	Author     []string
	Title      string
	Abstract   string
	Journal    string
	Date, Year string
	RIS        string
}

var pubref = map[string]string{}

// Register an author last name for cross referencing.
// E.g.:
// 	PubXRefAuthor("Vansteenkiste", "/people/arne")
// 	...
// 	LoadPublications("publications")
// will hyperlink all authors matching the regexp pattern "Vansteenkiste" to "/people/arne"
// It should be called before LoadPublications.
func PubXRefAuthor(regexp, href string) {
	pubref[regexp] = href
}

// Load publication .ciw files from the directory and serve them under that directory name.
// To be called before LoadContent. TODO: should be OK to call after loadcontent.
// To be called after PubXRefAuthor calls, if any.
func LoadPublications(dir string) {

	SetHandler("/"+dir, pubHandler)
	pubTemplate = loadTemplate(dir + "/template.html")

	ls := readDir(dir)
	for _, f := range ls {
		fullname := path.Clean(dir + "/" + f.Name())
		if path.Ext(f.Name()) == ".ciw" {
			pubs = append(pubs, parseRIS(fullname))
		}
	}

	pubcontent = newContent("Publication list", "", "")
	for i := range pubs {
		pubcontent.addChild(newContent("", "", pubs[i].Render()))
	}
}

func pubHandler(w http.ResponseWriter, r *http.Request) {
	renderContent(pubcontent, w, r)
}

func parseRIS(fname string) *pub {
	Log("parse", fname)
	p := new(pub)
	p.RIS = fname
	f, err := os.Open(fname)
	fatalErr(err)
	in := bufio.NewReader(f)

	l, _, err := in.ReadLine()
	key, val := string(l[:2]), string(l[3:])
	for len(l) > 3 {

		p.Add(key, val)

		l, _, err = in.ReadLine()
		k := string(l[:2])
		if k != "  " { // keep previous key if empty
			key = k
		}
		if len(l) > 2 {
			val = string(l[3:])
		}
	}
	//Log(p)
	return p
}

func (p *pub) Add(key, val string) {
	//	Log("pub", key, val)

	switch key {
	case "AF":
		val = xrefAuthor(val)
		p.Author = append(p.Author, val)
	case "TI":
		p.Title = p.Title + " " + val
	case "AB":
		p.Abstract = p.Abstract + " " + val
	case "JI":
		p.Journal = val
	case "PD":
		p.Date = val
	case "PY":
		p.Year = val
	}
}

// if name matches a key set by PubXRefAuthor, replace name by a hyperref to the corresponding link.
func xrefAuthor(name string) string {
	for k, v := range pubref {
		ok, _ := regexp.Match(k, []byte(name))
		if ok {
			return fmt.Sprintf(`<a href="%v" title="Go to personal page">%v</a>`, v, name)
		}
	}
	return name
}

func (p *pub) Render() string {
	var w bytes.Buffer
	fatalErr(pubTemplate.Execute(&w, p)) // TODO: not fatal, just panic and catch
	return w.String()
}
