package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth/gothic"
)

// func (cfg *apiConfig) handlerAuthCallbackFunction(w http.ResponseWriter, r *http.Request) {
// 	provider := chi.URLParam(r, "provider")
// 	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
// 	// fmt.Print(w)
// 	// fmt.Print("\n\n")
// 	// fmt.Print(r)
// 	// fmt.Print("\n\n")

// 	gothUser, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		RespondWithError(w, 400, fmt.Sprintf("can't complete user authentication: %v", err))
// 		return
// 	}
// 	fmt.Println("\nwe get here")
// 	http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
// 	cfg.handleUser(w, r, gothUser)
// }

func (cfg *apiConfig) handlerBeginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (cfg *apiConfig) handlerAuthLogout(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.Logout(w, r)
	w.Header().Set("Location", "http://localhost:5173")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
