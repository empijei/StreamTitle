package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/google/safehtml/template"

	"github.com/empijei/streamtitle/storage"
)

var tpls = template.Must(template.ParseGlob("templates/*.tpl.html"))

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	html, err := tpls.ExecuteTemplateToHTML(name, data)
	if err != nil {
		log.Printf("executing %q: %v", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	html, err = tpls.ExecuteTemplateToHTML("base.tpl.html", html)
	if err != nil {
		log.Printf("executing base.tpl.html: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, html.String())
}
func renderData(w http.ResponseWriter, data interface{}) {
	var buf bytes.Buffer
	if err := tpls.ExecuteTemplate(&buf, "base.tpl.html", data); err != nil {
		log.Printf("executing base.tpl.html: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

type noteTplData struct {
	Notes       []storage.Note
	OwnerOfPage string
	User        string
}
