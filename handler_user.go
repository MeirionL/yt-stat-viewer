package main

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handleUser(w http.ResponseWriter, r *http.Request, accessTok, refreshTok string) (uuid.UUID, error) {
	channelID, channelName, err := getAccountYTChannel(accessTok)
	if err != nil {
		return uuid.UUID{}, err
	}

	user, err := cfg.DB.GetUserByChannelID(r.Context(), channelID)

	if errors.Is(err, sql.ErrNoRows) {
		newUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:           uuid.New(),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			ChannelID:    channelID,
			ChannelName:  strings.ToLower(channelName),
			AccessToken:  accessTok,
			RefreshToken: refreshTok,
		})
		if err != nil {
			return uuid.UUID{}, err
		}

		return newUser.ID, nil
	}

	if err != nil {
		return uuid.UUID{}, err
	}

	err = cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:           user.ID,
		UpdatedAt:    time.Now().UTC(),
		ChannelName:  strings.ToLower(channelName),
		AccessToken:  accessTok,
		RefreshToken: refreshTok,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}
