package main

import "net/http"

const sessionCookie = "SESSION"

func setSession(w http.ResponseWriter, value string) {
	c := http.Cookie{Name: sessionCookie, Value: value}
	http.SetCookie(w, &c)
}

func (s *server) getUser(r *http.Request) string {
	var visitingUser string
	sess, err := r.Cookie(sessionCookie)
	if err == nil && sess.Value != "" {
		visitingUser, _ = s.sessions.GetName(sess.Value)
	}
	return visitingUser
}

func (s *server) checkAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usr := s.getUser(r)
		if usr == "" {
			renderPage(w, `Please <a href="/">login</a> first.`)
			return
		}
		h.ServeHTTP(w, r)
	}
}
