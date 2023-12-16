package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

func (cfg *apiConfig) handlerAuthLogin(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	// fmt.Print(w)
	// fmt.Print("\n\n")
	// fmt.Print(r)
	// fmt.Print("\n\n")

	if user, err := gothic.CompleteUserAuth(w, r); err == nil {
		fmt.Println("we get here")
		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
		cfg.handleUser(w, r, user)
		return
	} else {
		fmt.Println("it is time to begin logging in!")
		gothic.BeginAuthHandler(w, r)
	}
}

func (cfg *apiConfig) handlerAuthLogout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "http://localhost:5173")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
