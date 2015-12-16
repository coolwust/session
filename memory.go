package session

import (
	"sync"
	"time"
)

var _ Store = new(MemoryStore)

type MemoryStore struct {
	data map[string]*Session
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]*Session),
	}
}

func (store *MemoryStore) Set(session *Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[session.ID] = session
	return nil
}

func (store *MemoryStore) Get(id string) (*Session, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	session, ok := store.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return session, nil
}

func (store *MemoryStore) Remove(id string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.data, id)
	return nil
}

func (store *MemoryStore) Clean() error {
	store.mu.Lock()
	defer store.mu.Unlock()
	now := time.Now()
	for id, session := range store.data {
		if session.Expires.Before(now) {
			delete(store.data, id)
		}
	}
	return nil
}

var _ Store = &MemoryStore{}
