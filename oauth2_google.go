package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var conf *oauth2.Config

func google_auth(w http.ResponseWriter, r *http.Request) {
	fmt.Print("\n at least we started\n")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	// Your credentials should be obtained from the Google
	// Developer Console (https://console.developers.google.com).
	conf = &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/youtube.channel-memberships.creator",
			"https://www.googleapis.com/auth/yt-analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}
	// Redirect user to Google's consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline) // asks for refresh token and specific scope permissions
	fmt.Printf("url is: %v", url)
	http.Redirect(w, r, url, http.StatusFound)
}

func (cfg *apiConfig) get_google_token(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code parameter missing", http.StatusBadRequest)
		return
	}

	// Handle the exchange code to initiate a transport.
	tok, err := conf.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nrefresh token is: %v\n", tok.RefreshToken)
	// // code to refresh access token
	// if tok.Expiry.Before(time.Now()) {
	// 	// If the token is expired, use the refresh token to obtain a new token
	// 	tok, err = conf.TokenSource(context.Background(), tok).Token()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	fmt.Print("\nwe made it here!!!!!\n")

	userID := cfg.handleUser(w, r, tok.AccessToken, tok.RefreshToken)

	// Redirect first, then set the message in a cookie or as a query parameter
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
		fmt.Printf("no user to logout with id of %v", channelID)
		return
	} else if err != nil {
		fmt.Printf("error deleting user with id: %v.\n Error: %v", channelID, err)
		return
	} else {
		fmt.Printf("logged out account: %v", deletedChannelName)
		return
	}
}
