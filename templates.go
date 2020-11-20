package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/empijei/streamtitle/storage"
)

var tpls = template.Must(template.ParseGlob("templates/*.tpl.html"))

var _ template.HTML

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var buf bytes.Buffer
	if err := tpls.ExecuteTemplate(&buf, name, data); err != nil {
		log.Printf("executing %q: %v", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	str := buf.String()
	html := template.HTML(str)
	buf.Reset()
	if err := tpls.ExecuteTemplate(&buf, "base.tpl.html", html); err != nil {
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
	Notes         []storage.Note
	OwnerOfPage   string
	User          string
	OwnerOfPageJS string
	UserJS        string
}

func (n *noteTplData) escape() {
	n.OwnerOfPage = template.HTMLEscapeString(n.OwnerOfPage)
	n.User = template.HTMLEscapeString(n.User)
	n.OwnerOfPage = template.HTMLEscapeString(template.JSEscapeString(n.OwnerOfPage))
	n.User = template.HTMLEscapeString(template.JSEscapeString(n.User))
	for i, nb := range n.Notes {
		nb.Text = template.HTMLEscapeString(nb.Text)
		nb.Title = template.HTMLEscapeString(nb.Title)
		n.Notes[i] = nb
	}
}
