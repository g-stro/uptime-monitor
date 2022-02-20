package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func index(writer http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(writer, "index.html", record)
}
