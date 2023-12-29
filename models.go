package main

import (
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ChannelID    string    `json:"channel_id"`
	ChannelName  string    `json:"channel_name"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

func databaseUserToUser(dbUser database.User) User {
	return User{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		ChannelID:    dbUser.ChannelID,
		ChannelName:  dbUser.ChannelName,
		AccessToken:  dbUser.AccessToken,
		RefreshToken: dbUser.RefreshToken,
	}
}

func databaseUsersToUsers(dbUsers []database.User) []User {
	users := []User{}
	for _, dbUser := range dbUsers {
		users = append(users, databaseUserToUser(dbUser))
	}
	return users
}
