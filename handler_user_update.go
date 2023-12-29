package main

import (
	"fmt"
	"net/http"

	"github.com/MeirionL/stream_stats/backend/internal/database"
)

func (cfg *apiConfig) handlerDeleteUser(w http.ResponseWriter, r *http.Request, user database.User) {

	_, err := cfg.DB.DeleteUser(r.Context(), user.ID)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("couldn't delete user: %v", err))
		return
	}
	RespondWithJSON(w, 200, struct{}{})
}
