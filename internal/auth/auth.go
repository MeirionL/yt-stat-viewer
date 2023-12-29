package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitch"
)

type GoogleClaims struct {
	Email string `json:"email"`

	EmailVerified bool `json:"email_verified"`

	FirstName string `json:"given_name"`

	LastName string `json:"family_name"`

	jwt.StandardClaims
}

func ValidateGoogleJWT(tokenString string) (GoogleClaims, error) {
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != "1047286383284-a25hpilnspp1ttpe2the8ml5juaogbsd.apps.googleusercontent.com" {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}

func getGooglePublicKey(keyID string) (string, error) {

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")

	if err != nil {
		return "", err
	}

	dat, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	myResp := map[string]string{}

	err = json.Unmarshal(dat, &myResp)

	if err != nil {
		return "", err
	}

	key, ok := myResp[keyID]

	if !ok {
		return "", errors.New("key not found")
	}

	return key, nil
}

func MakeJWT(email, jwtSecret string) (string, error) {
	signingKey := []byte(jwtSecret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   fmt.Sprintf("%v", email),
	})

	return token.SignedString(signingKey)
}

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
