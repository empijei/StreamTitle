package storage

import (
	"errors"
	"sync"
)

// Credentials is a storage for credentials.
// There are several security issues here, please don't use this in production.
type Credentials struct {
	mu    sync.Mutex
	users map[string]string
}

// NewCredentials creates a new credentials storage.
func NewCredentials() *Credentials {
	// TODO savefile
	return &Credentials{users: map[string]string{}}
}

// HasUser checks if the user exists.
func (c *Credentials) HasUser(name string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, has := c.users[name]
	return has
}

// ErrAlreadyExists is returned when attempting to store something that already is in store.
var ErrAlreadyExists = errors.New("Already exists")

// AddUser adds a user to the storage if it is not already there.
func (c *Credentials) AddUser(name, password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, has := c.users[name]; has {
		return ErrAlreadyExists
	}
	c.users[name] = password
	return nil
}

// AuthUser validates the provided credentials against the storage.
func (c *Credentials) AuthUser(name, password string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	pw, has := c.users[name]
	if !has {
		return false
	}
	return pw == password
}
