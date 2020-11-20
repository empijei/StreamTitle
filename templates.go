package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/empijei/streamtitle/storage"
)

var tpls = template.Must(template.ParseGlob("templates/*.tpl.html"))

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var buf bytes.Buffer
	if err := tpls.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("executing %q: %v", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	str := buf.String()
	buf.Reset()
	if err := tpls.ExecuteTemplate(&buf, "base.tpl.html", str); err != nil {
		log.Printf("executing base.tpl.html: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
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

func (n *noteTplData) escape() {
	n.OwnerOfPage = template.HTMLEscapeString(n.OwnerOfPage)
	n.User = template.HTMLEscapeString(n.User)
	for i, nb := range n.Notes {
		nb.Text = template.HTMLEscapeString(nb.Text)
		nb.Title = template.HTMLEscapeString(nb.Title)
		n.Notes[i] = nb
	}
}
