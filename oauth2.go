package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func google_auth(w http.ResponseWriter, r *http.Request) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	conf = &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/youtube.readonly",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	http.Redirect(w, r, url, http.StatusFound)
}

func (cfg *apiConfig) get_google_token(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		RespondWithError(w, 400, "code parameter missing from request url")
		return
	}

	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(r.Context(), code)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("couldn't exchange code for token: %v", err))
		return
	}

	userID, err := cfg.handleUser(w, r, tok.AccessToken, tok.RefreshToken)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error handling user: %v", err))
		return
	}

	// Redirect first, then set the needed channel id as a query parameter
	http.Redirect(w, r, "http://localhost:5173?id="+url.QueryEscape(userID.String()), http.StatusFound)
}

func (cfg *apiConfig) logout(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "id")
	channelUUID, err := uuid.Parse(channelID)
	if err != nil {
		RespondWithError(w, 400, "couldn't parse ID to UUID value")
	}

	deletedChannelName, err := cfg.DB.DeleteUser(r.Context(), channelUUID)
	if errors.Is(err, sql.ErrNoRows) {
		RespondWithError(w, 400, fmt.Sprintf("no user to logout with id of %v", channelID))
		return
	} else if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error deleting user with id: %v", channelID))
		return
	} else {
		fmt.Printf("logged out user with channel name: %v\n", deletedChannelName)
		http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
		return
	}
}

func (cfg *apiConfig) refreshAccessToken(u database.User) (database.User, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/callback",
	}

	token := &oauth2.Token{
		RefreshToken: u.RefreshToken,
	}

	newToken, err := config.TokenSource(ctx, token).Token()
	if err != nil {
		return database.User{}, err
	}

	if newToken.AccessToken != u.AccessToken && newToken.RefreshToken == u.RefreshToken {
		newUser, err := cfg.DB.UpdateAccessToken(ctx, database.UpdateAccessTokenParams{
			ID:          u.ID,
			UpdatedAt:   time.Now().UTC(),
			AccessToken: newToken.AccessToken,
		})
		if err != nil {
			return database.User{}, err
		}

		return newUser, nil
	} else {
		return database.User{}, fmt.Errorf("discrepency in refrehed token. Can't authorize for channel: %v", u.ChannelName)
	}
}
