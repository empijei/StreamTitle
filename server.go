package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/empijei/streamtitle/storage"
)

type server struct {
	mux      *http.ServeMux
	notes    *storage.Notes
	sessions *storage.Sessions
	creds    *storage.Credentials
}

func newServer() *server {
	s := server{
		mux:      http.NewServeMux(),
		notes:    storage.NewNotes(),
		sessions: storage.NewSessions(),
		creds:    storage.NewCredentials(),
	}
	s.mux.Handle("/notes/", s.notesHandler())
	s.mux.Handle("/me", s.meHandler())
	s.mux.Handle("/login", s.loginHandler())
	s.mux.Handle("/addnote", s.addNoteHandler())
	s.mux.Handle("/static/", http.StripPrefix("/static/", s.staticHandler()))
	s.mux.Handle("/", s.indexHandler())
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
		usr := s.getUser(r)
		if usr == "" {
			renderPage(w, loginForm)
			return
		}
		http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
	}
}

func (s *server) notesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		visitingUser := s.getUser(r)
		pos := strings.LastIndex(r.URL.Path, "/")
		viewedUser := r.URL.Path[pos+1:]
		var notes []storage.Note
		if visitingUser == viewedUser {
			// Render all notes, the user is viewing its own data.
			notes = s.notes.GetNotes(viewedUser, true /*show private*/)
		} else {
			// Only render public notes.
			notes = s.notes.GetNotes(viewedUser, false /*no show private*/)
		}
		var sb strings.Builder
		data := map[string]interface{}{
			"Notes":        notes,
			"ViewedUser":   viewedUser,
			"VisitingUser": visitingUser,
		}
		if err := listNotesTpl.Execute(&sb, data); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		renderPage(w, sb.String())
	}
}

func (s *server) loginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		uname := r.FormValue("uname")
		passwd := r.FormValue("password")
		if s.creds.HasUser(uname) {
			if !s.creds.AuthUser(uname, passwd) {
				renderPage(w, `Invalid password, please <a href="/">retry logging in.</a>`)
				return
			}
		} else {
			s.creds.AddUser(uname, passwd)
		}
		setSession(w, s.sessions.GetToken(uname))
		http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
	}
}

func (s *server) meHandler() http.HandlerFunc {
	return s.checkAuth(func(w http.ResponseWriter, r *http.Request) {
		usr := s.getUser(r)
		notes := s.notes.GetNotes(usr, true /*show private*/)
		var sb strings.Builder
		data := map[string]interface{}{
			"Notes": notes,
			"User":  usr,
		}
		if err := myNotesTpl.Execute(&sb, data); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		renderPage(w, sb.String())
	})
}

func (s *server) addNoteHandler() http.HandlerFunc {
	return s.checkAuth(func(w http.ResponseWriter, r *http.Request) {
		s.notes.AddNote(s.getUser(r), storage.Note{
			Title:   r.FormValue("title"),
			Text:    r.FormValue("text"),
			Private: r.FormValue("private") != "",
		})
		http.Redirect(w, r, "/me", http.StatusTemporaryRedirect)
	})
}
