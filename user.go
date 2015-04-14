package main

import "golang.org/x/oauth2"

func NewUser(idToken *googleIdToken, authToken *oauth2.Token) *User {
	return &User{
		idToken:     idToken,
		googleToken: authToken,
	}
}

type User struct {
	idToken     *googleIdToken
	profile     *googleProfile
	googleToken *oauth2.Token
}

func (u *User) Email() string {
	return u.idToken.Email
}

func (u *User) String() string {
	if u.profile != nil {
		return u.profile.DisplayName
	}
	return u.idToken.Email
}

func (u *User) SetProfile(profile *googleProfile) {
	u.profile = profile
}
