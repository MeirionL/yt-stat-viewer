package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request, accessTok, refreshTok string) uuid.UUID {
	channelID, channelName, err := getAccountYTChannel(accessTok)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("can't get channel name: %v", err))
		return uuid.UUID{}
	}

	user, err := cfg.DB.GetUserByChannelID(r.Context(), channelID)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("\n%v is a new channel, so we are going to make a user for them!\n", channelName)
		newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			ChannelID:    channelID,
			ChannelName:  channelName,
			AccessToken:  accessTok,
			RefreshToken: refreshTok,
		})
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("couldn't create user: %v", err))
			return uuid.UUID{}
		}

		fmt.Println(newUser)
		return newUser.ID
	} else if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("can't get user by channel id: %v", err))
		return uuid.UUID{}
	}

	err = cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           user.ID,
		UpdatedAt:    time.Now().UTC(),
		ChannelName:  channelName,
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
	})
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("couldn't update tokens: %v", err))
		return uuid.UUID{}
	}

	return user.ID
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		ChannelID    string `json:"channel_id"`
		ChannelName  string `json:"channel_name"`
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
		ChannelID:    params.ChannelID,
		ChannelName:  params.ChannelName,
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
