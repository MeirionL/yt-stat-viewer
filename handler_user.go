package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/google/uuid"
	"github.com/markbates/goth"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request, u goth.User) {
	var platform string
	if u.Provider == "google" {
		platform = "youtube"
	} else if u.Provider == "twitch" {
		platform = "twitch"
	} else {
		RespondWithError(w, 400, "Invalid platform")
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Email:        u.Email,
		Platform:     platform,
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error creating user: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(databaseUserToUser(user))
	if err != nil {
		log.Printf("error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Write(dat)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email        string `json:"email"`
		Platform     string `json:"platform"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error parsing JSON: %v", err))
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Email:        params.Email,
		Platform:     params.Platform,
		AccessToken:  params.AccessToken,
		RefreshToken: params.RefreshToken,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("error creating user: %v", err))
		return
	}

	RespondWithJSON(w, 201, databaseUserToUser(user))
}

func (cfg *apiConfig) handlerGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := cfg.DB.GetUsers(r.Context())
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("couldn't get users: %v", err))
		return
	}

	RespondWithJSON(w, 200, databaseUsersToUsers(users))
}

func (cfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJSON(w, 200, databaseUserToUser(user))
}
