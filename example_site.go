// +build ignore

package main

import (
	. "github.com/barnex/webserv"
	"os"
	"net/http"
)

func main() {

	SetHandler("/uname", Command("uname", "-a"))
	SetHandler("/fortune", Command("fortune"))
	SetHandler("/top", Command("top", "-b", "-n", "1"))

	SetHandler("/restart", func(http.ResponseWriter, *http.Request){
		Log("restart")
		os.Exit(0)
	})

	PubXRefAuthor("Waeyenberge", "/people/bartel")
	PubXRefAuthor("Vansteenkiste", "/people/arne")
	PubXRefAuthor("Dvornik", "/people/mykola")
	PubXRefAuthor("Helsen", "/people/mathias")
	LoadPublications("publications")

	LoadContent(".")

	Run()
}
