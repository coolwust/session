package session

import (
	"sync"
	"time"
	"crypto/rand"
	"fmt"
	"net/http"
)

type Session struct {
	ID       string
	Expires  time.Time
	Duration time.Duration
	data     map[string]interface{}
	mu       sync.RWMutex
}

func NewSession(genID func() string, duration time.Duration) *Session {
	session := &Session{
		ID:       genID(),
		Expires:  time.Now().Add(duration),
		Duration: duration,
		data:     make(map[string]interface{}),
	}
	return session
}

func (session *Session) Set(name string, value interface{}) {
	session.mu.Lock()
	defer session.mu.Unlock()
	session.data[name] = value
}

func (session *Session) Get(name string) interface{} {
	session.mu.RLock()
	defer session.mu.RUnlock()
	return session.data[name]
}

func (session *Session) Touch() {
	session.Expires = time.Now().Add(session.Duration)
}

func FromRequest(req *http.Request, name, key string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}
	sid, err := Unsign(cookie.Value, key)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func ToResponse(w http.ResponseWriter, sid, name, path string, maxAge int, key string) {
	cookie := &http.Cookie{
		Name:   name,
		Path:   path,
		MaxAge: maxAge,
	}
	if maxAge != -1 {
		cookie.Value = Sign(sid, key)
		//cookie.Secure = true
		//cookie.HttpOnly = true
	}
	w.Header().Set("Set-Cookie", cookie.String())
}

func UUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:])
}
