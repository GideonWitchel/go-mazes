package main

import (
	"html/template"
	"net/http"
)

var tpl = template.Must(template.ParseFiles("templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	port := "3000"
	http.ListenAndServe(":"+port, mux)
}
