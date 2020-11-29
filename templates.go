package main

import (
	"io"
	"log"
	"net/http"

	"github.com/google/safehtml/template"

	"github.com/empijei/streamtitle/csp"
	"github.com/empijei/streamtitle/storage"
)

var tpls = template.Must(template.ParseGlob("templates/*.tpl.html"))

func renderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	html, err := tpls.ExecuteTemplateToHTML(name, data)
	if err != nil {
		log.Printf("executing %q: %v", name, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	html, err = tpls.ExecuteTemplateToHTML("base.tpl.html", baseTplData{
		CSPNonce: csp.GetNonce(r),
		Data:     html,
	})
	if err != nil {
		log.Printf("executing base.tpl.html: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, html.String())
}

func renderData(w http.ResponseWriter, r *http.Request, data interface{}) {
	html, err := tpls.ExecuteTemplateToHTML("base.tpl.html", baseTplData{
		CSPNonce: csp.GetNonce(r),
		Data:     data,
	})
	if err != nil {
		log.Printf("executing base.tpl.html: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, html.String())
}

type baseTplData struct {
	CSPNonce string
	Data     interface{}
}

type noteTplData struct {
	Notes       []storage.Note
	OwnerOfPage string
	User        string
}
