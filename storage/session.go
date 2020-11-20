package storage

import (
	"encoding/base64"
	"math/rand"
	"sync"
)

// Sessions is a session storage.
// There are several security issues here, please don't use this in production.
type Sessions struct {
	mu    sync.RWMutex
	byTok map[string]string
	byUsr map[string]string
}

// NewSessions creates a new sessions storage
func NewSessions() *Sessions {
	return &Sessions{
		byTok: map[string]string{},
		byUsr: map[string]string{},
	}
}

// GetName returns the username for the provided token if available.
func (s *Sessions) GetName(token string) (user string, valid bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, valid = s.byTok[token]
	return
}

// GetToken returns the session token for the given user.
// If not available one will be generated.
func (s *Sessions) GetToken(user string) (token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	token, has := s.byUsr[user]
	if has {
		return token
	}
	token = genToken()
	s.byUsr[user] = token
	s.byTok[token] = user
	return token
}

// DelTokenForUser removes the token for the given user.
func (s *Sessions) DelTokenForUser(user string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	token, has := s.byUsr[user]
	if !has {
		return
	}
	delete(s.byTok, token)
	delete(s.byUsr, user)
}

func genToken() string {
	b := make([]byte, 20)
	rand.Read(b)
	tok := base64.RawStdEncoding.EncodeToString(b)
	return tok
}
