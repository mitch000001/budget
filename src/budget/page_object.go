package main

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

func PageForSession(s *Session) *PageObject {
	p := make(PageObject)
	p["session"] = s
	p["RequestPath"] = s.URL.Path
	if s.User != nil {
		p["user"] = s.User
	}
	sessErrors := s.GetErrors()
	if len(sessErrors) > 0 {
		p.AddErrors(sessErrors)
		s.ResetErrors()
	}
	return &p
}

func NewPageObject() *PageObject {
	p := make(PageObject)
	return &p
}

type PageObject map[string]interface{}

func (p *PageObject) LoggedIn() bool {
	s, ok := (*p)["session"]
	if !ok {
		return false
	}
	return s.(*Session).LoggedIn()
}

func (p *PageObject) Debug() bool {
	return debugMode
}

func (p *PageObject) CurrentUser() *User {
	return (*p)["user"].(*User)
}

func (p *PageObject) Errors() []error {
	errs, ok := (*p)["errors"]
	if !ok {
		return nil
	} else {
		return errs.([]error)
	}
}

func (p *PageObject) AddError(err error) {
	errs, ok := (*p)["errors"]
	var errors []error
	if !ok || errs == nil {
		errors = make([]error, 0)
	} else {
		errors = errs.([]error)
	}
	errors = append(errors, err)
	(*p)["errors"] = errors
}

func (p *PageObject) AddErrors(errs []error) {
	pErrs, ok := (*p)["errors"]
	var errors []error
	if !ok || errs == nil {
		errors = make([]error, 0)
	} else {
		errors = pErrs.([]error)
	}
	errors = append(errors, errs...)
	(*p)["errors"] = errors
}

func (p *PageObject) Set(key string, value interface{}) {
	if isLower(key) {
		panic(fmt.Errorf("Key must begin with a capital letter"))
	}
	(*p)[key] = value
}

func isLower(input string) bool {
	c, _ := utf8.DecodeRuneInString(input)
	return unicode.IsLower(c)
}
