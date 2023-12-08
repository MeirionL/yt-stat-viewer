package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name       string `json:"name"`
		OAuthToken string `json:"oauth_token"`
	}
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't decode parameters: %v", err))
		return
	}

	newUser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:         user.ID,
		UpdatedAt:  time.Now(),
		Name:       params.Name,
		OauthToken: params.OAuthToken,
	})
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't update user: %v", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        newUser.ID,
			UpdatedAt: newUser.UpdatedAt,
			Name:      newUser.Name,
			APIKey:    newUser.ApiKey,
		},
	})
}

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {

	err := cfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("couldn't delete user: %v", err))
		return
	}
	RespondWithJSON(w, 200, struct{}{})
}
