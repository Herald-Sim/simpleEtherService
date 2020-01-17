package main

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("./static/index.html"))

// for security issue checking
var validPath = regexp.MustCompile("^/(index|regist|save)/")

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index")
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/index/", indexHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

	//log.Fatal(http.ListenAndServeTLS(":443", "/etc/letsencrypt/live/lbankkorea.com/fullchain.pem", "/etc/letsencrypt/live/lbankkorea.com/privkey.pem", nil))
}
