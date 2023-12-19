package auth

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitch"
)

const (
	MaxAge = 86400 * 30
	IsProd = false
)

func NewAuth() {
	godotenv.Load(".env")

	key := os.Getenv("RANDOM_KEY")
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	twitchClientSecret := os.Getenv("TWITCH_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(
		google.New(
			googleClientID, googleClientSecret, "http://localhost:8080/auth/google/callback",
			"https://www.googleapis.com/auth/youtube.readonly",
			"https://www.googleapis.com/auth/youtube.channel-memberships.creator",
			"https://www.googleapis.com/auth/yt-analytics.readonly",
		),
		twitch.New(
			twitchClientID, twitchClientSecret, "http://localhost:8080/auth/twitch/callback",
			"channel:read:subscriptions",
			"moderation:read",
			"moderator:read:followers",
			"user:read:broadcast",
			"user:read:email",
		),
	)
}

// // Example:
// // Authorisation: ApiKey {insert apikey here}
// func GetAPIKey(headers http.Header) (string, error) {
// 	val := headers.Get("Authorization")
// 	if val == "" {
// 		return "", errors.New("no authentication info found")
// 	}

// 	vals := strings.Split(val, " ")
// 	if len(vals) != 2 {
// 		return "", errors.New("malformed auth header")
// 	}
// 	if vals[0] != "ApiKey" {
// 		return "", errors.New("malformed first part of auth header")
// 	}
// 	return vals[1], nil
// }
