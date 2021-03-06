package main

import (
	"net/http"
	"strings"

	"github.com/empijei/streamtitle/csp"
	"github.com/empijei/streamtitle/storage"
)

type server struct {
	mux      http.Handler
	notes    *storage.Notes
	sessions *storage.Sessions
	creds    *storage.Credentials
}

func newServer() *server {
	s := server{
		notes:    storage.NewNotes(),
		sessions: storage.NewSessions(),
		creds:    storage.NewCredentials(),
	}
	mux := http.NewServeMux()
	mux.Handle("/notes/", s.notesHandler())
	mux.Handle("/login", s.loginHandler())
	mux.Handle("/logout", s.logoutHandler())
	mux.Handle("/addnote", s.addNoteHandler())
	mux.Handle("/static/", http.StripPrefix("/static/", s.staticHandler()))
	mux.Handle("/", s.indexHandler())
	s.mux = csp.Protect(mux, true)
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) staticHandler() http.HandlerFunc {
	dir := http.Dir("static")
	fs := http.FileServer(dir)
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := dir.Open(r.URL.Path)
		if err != nil {
			// Do not reveal information on the error.
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		fs.ServeHTTP(w, r)
	}
}

func (s *server) indexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if p := r.URL.Path; p != "" && p != "/index.html" && p != "/" {
			http.Redirect(w, r, "/notes/"+r.URL.Path, http.StatusTemporaryRedirect)
			return
		}
		usr := s.getUser(r)
		if usr == "" {
			renderTemplate(w, r, "login.tpl.html", nil)
			return
		}
		http.Redirect(w, r, "/notes/"+usr, http.StatusTemporaryRedirect)
	}
}

func (s *server) notesHandler() http.HandlerFunc {
	return s.checkAuth(func(w http.ResponseWriter, r *http.Request) {
		user := s.getUser(r)
		ownerOfPage := r.FormValue("user")
		if ownerOfPage == "" {
			pos := strings.LastIndex(r.URL.Path, "/")
			ownerOfPage = r.URL.Path[pos+1:]
		}
		var notes []storage.Note
		if user == ownerOfPage {
			// Render all notes, the user is viewing its own data.
			notes = s.notes.GetNotes(ownerOfPage, true /*show private*/)
		} else {
			// Only render public notes.
			notes = s.notes.GetNotes(ownerOfPage, false /*no show private*/)
		}
		data := noteTplData{
			Notes:       notes,
			OwnerOfPage: ownerOfPage,
			User:        user,
		}
		//data.escape()
		renderTemplate(w, r, "notes.tpl.html", data)
	})
}

func (s *server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uname := r.FormValue("uname")
		passwd := r.FormValue("password")
		if s.creds.HasUser(uname) {
			if !s.creds.AuthUser(uname, passwd) {
				renderTemplate(w, r, "invalid_pw.tpl.html", uname)
				return
			}
		} else {
			s.creds.AddUser(uname, passwd)
		}
		setSession(w, s.sessions.GetToken(uname))
		http.Redirect(w, r, "/notes/"+uname, http.StatusTemporaryRedirect)
	}
}

func (s *server) logoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setSession(w, "")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func (s *server) addNoteHandler() http.HandlerFunc {
	return s.checkAuth(func(w http.ResponseWriter, r *http.Request) {
		usr := r.FormValue("target_user")
		n := storage.Note{
			Title: r.FormValue("title"),
			Text:  r.FormValue("text"),
		}
		if usr == s.getUser(r) { // Only allow owners to set private bit
			n.Private = r.FormValue("private") != ""
		}
		s.notes.AddNote(usr, n)
		http.Redirect(w, r, "/notes/"+usr, http.StatusTemporaryRedirect)
	})
}
