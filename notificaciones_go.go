package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

var calls = 0

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}
type Data struct {
	Calls int
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	calls = calls + 1

	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	data := Data{
		Calls: calls,
	}

	t.templ.Execute(w, &data)
}

func testHello(w http.ResponseWriter, r *http.Request) {
	calls = calls + 1;
	fmt.Fprintf(w, "Hello, there! For the %d time.", calls);
}

func testHelloHtml(w http.ResponseWriter, r *http.Request) {
	calls = calls + 1;

	fmt.Fprintf(w, "<h1 class='ui header'>" +
	"Welcome to go!" +
	"<div class='sub header'>For the %d time!</div>" +
	"</h1>",
		calls)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse() // parse the flags

	http.Handle("/", &templateHandler{filename: "welcome.html"})
	http.HandleFunc("/html", testHelloHtml)
	http.HandleFunc("/test", testHello)

	// start the web server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}