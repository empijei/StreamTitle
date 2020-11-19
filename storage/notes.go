package storage

import (
	"sort"
	"sync"
)

// Note is a note.
type Note struct {
	Title, Text string
	Private     bool
}

// Notes is a notes storage.
// There are several security issues here, please don't use this in production.
type Notes struct {
	mu    sync.RWMutex
	notes map[string]map[string]Note
}

// NewNotes creates a new notes storage
func NewNotes() *Notes {
	return &Notes{notes: map[string]map[string]Note{}}
}

// GetNotes retrieves all the notes for the given user, alphabetically ordered.
func (n *Notes) GetNotes(user string, private bool) []Note {
	n.mu.RLock()
	defer n.mu.RUnlock()
	notes := n.notes[user]
	var ns []Note
	for _, n := range notes {
		if private || n.Private == false {
			ns = append(ns, n)
		}
	}
	sort.Slice(ns, func(i, j int) bool { return ns[i].Title < ns[j].Title })
	return ns
}

// AddNote adds or overwrite the given note for the given user.
func (n *Notes) AddNote(user string, add Note) {
	n.mu.Lock()
	defer n.mu.Unlock()
	u := n.notes[user]
	if u == nil {
		u = map[string]Note{}
	}
	u[add.Title] = add
	n.notes[user] = u
}

/*
// GetNote retrieves the given note for the given user.
// If the note was not found the second return value is false.
func (n *Notes) GetNote(user, title string) (note Note, has bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	note, has = n.notes[user][title]
	return
}
*/
