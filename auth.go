package main

import (
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var Users = map[string]string{
	"admin": "admin",
	"mike":  "password1",
	"john":  "123456",
}

var sessions = map[string]session{}

type session struct {
	username string
	expiry   time.Time
}

func (s session) isExpired() bool {
	fmt.Println("Expiry:", s.expiry, time.Now())
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleSignInSubmit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// Handle the error
		http.Error(w, "Unable to parse form data", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	creds := Credentials{username, password}

	fmt.Println("Username:", creds.Username, "Password:", creds.Password)

	expectedPassword, ok := Users[creds.Username]
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Println("Authentication failed for user:", creds.Username)
		return
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Hour)

	sessions[sessionToken] = session{creds.Username, expiresAt}

	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		Path:     "/",
		HttpOnly: false,
	}

	http.SetCookie(w, cookie)

	fmt.Println("User:", creds.Username, "logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetUsernameFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	sessionToken := cookie.Value

	userSession, ok := sessions[sessionToken]

	if userSession.isExpired() {
		fmt.Println("Session expired for user:", userSession.username)
		delete(sessions, sessionToken)
		return ""
	}

	if ok {
		fmt.Println("User is logged in", userSession.username)
		return userSession.username
	}

	fmt.Println("User is not logged in")
	return ""
}

func handleSignOut(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := cookie.Value

	delete(sessions, sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
