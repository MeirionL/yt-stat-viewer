package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) getKickChannelStats(w http.ResponseWriter, r *http.Request) {
	channelString := chi.URLParam(r, "channel")

	newChannel := Channel{
		Subscribers: 5,
		Title:       channelString,
		Views:       200,
		Platform:    "Kick",
	}
	cfg.channels = append(cfg.channels, newChannel)

	RespondWithJSON(w, http.StatusOK, cfg.channels)
}
