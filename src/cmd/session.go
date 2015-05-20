package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type SessionManager map[string]*Session

func (s *SessionManager) init() {
	if s == nil {
		*s = make(map[string]*Session)
	}
}

func (sm *SessionManager) Add(s *Session) {
	sm.init()
	(*sm)[s.id] = s
}

func (sm *SessionManager) Find(sessionId string) *Session {
	return (*sm)[sessionId]
}

func (sm *SessionManager) Remove(s *Session) {
	delete(*sm, s.id)
}

func GetSessionFromCookie(cookie *http.Cookie) *Session {
	if cookie == nil {
		return nil
	}
	if expired := cookie.Expires.After(time.Now()); expired {
		return nil
	}
	sessionId := cookie.Value
	return sessions.Find(sessionId)
}

type Session struct {
	Stack    string
	URL      *url.URL
	location string
	User     *User
	id       string
	errors   []error
}

func (s *Session) LoggedIn() bool {
	return s.User != nil
}

func (s *Session) AddError(err error) {
	if s.errors == nil {
		s.errors = make([]error, 0)
	}
	s.errors = append(s.errors, err)
}

func (s *Session) AddDebugError(err error) {
	if debugMode {
		s.AddError(err)
	}
}

func (s *Session) GetErrors() []error {
	return s.errors
}

func (s *Session) ResetErrors() {
	s.errors = make([]error, 0)
}

func newSession() *Session {
	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	id := fmt.Sprintf("%x", sha256.Sum256(b))
	return &Session{id: id}
}
