package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jws"
)

type googleIdToken struct {
	jws.ClaimSet
	AccessTokenHash     string `json:"at_hash"`
	EmailVerified       bool   `json:"email_verified"`
	AuthorizedPresenter string `json:"azp"`
	Email               string `json:"email"`
	HostedDomain        string `json:"hd"`
}

type googleProfile struct {
	DisplayName string `json:"displayName"`
	Name        struct {
		Formatted  string `json:"formatted"`
		FamilyName string `json:"familyName"`
		GivenName  string `json:"givenName"`
		MiddleName string `json:"middleName"`
	}
	Domain string `json:"domain"`
}

var profileUrl string = "https://www.googleapis.com/plus/v1/people/me"

func (g *googleProfile) FullName() string {
	return fmt.Sprintf("%s %s", g.Name.GivenName, g.Name.FamilyName)
}

func googleLoginHandler(config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := newSession()
		s.location = r.Header.Get("X-Referer")
		sessions.Add(s)
		url := config.AuthCodeURL(s.id, oauth2.AccessTypeOffline)
		http.Redirect(w, r, url, http.StatusFound)
		return
	}
}

func googleRedirectHandler(config *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		state := params.Get("state")
		if state == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		session := sessions.Find(state)
		if session == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		code := params.Get("code")
		if code == "" {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		token, err := config.Exchange(oauth2.NoContext, code)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		id := token.Extra("id_token")
		idToken, err := decode(id.(string))
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user := NewUser(idToken, token)
		session.User = user
		defer func() {
			http.SetCookie(w, NewSessionCookie(session.id))
			http.Redirect(w, r, session.location, http.StatusFound)
		}()
		client := config.Client(oauth2.NoContext, token)
		response, err := client.Get(profileUrl)
		if err != nil {
			log.Printf("Error %T: %v\n", err, err)
			return
		}
		defer response.Body.Close()
		profile, err := decodeGoogleProfile(response.Body)
		if err != nil {
			log.Printf("Error %T: %v\n", err, err)
			return
		}
		session.User.SetProfile(profile)
	}
}

func decodeGoogleProfile(r io.Reader) (*googleProfile, error) {
	var profile googleProfile
	err := json.NewDecoder(r).Decode(&profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func NewSessionCookie(sessionId string) *http.Cookie {
	return &http.Cookie{Name: "budget", Value: sessionId, Expires: time.Now().Add(5 * 24 * time.Hour)}
}

func decode(payload string) (*googleIdToken, error) {
	// decode returned id token to get expiry
	s := strings.Split(payload, ".")
	if len(s) < 2 {
		return nil, errors.New("jws: invalid token received")
	}
	decoded, err := base64Decode(s[1])
	if err != nil {
		return nil, err
	}
	c := &googleIdToken{}
	err = json.NewDecoder(bytes.NewBuffer(decoded)).Decode(c)
	return c, err
}

func base64Decode(s string) ([]byte, error) {
	// add back missing padding
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}
