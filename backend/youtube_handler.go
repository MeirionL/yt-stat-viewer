package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YoutubeStats struct {
	Subscribers    int    `json:"subscribers"`
	ChannelName    string `json:"channel_name"`
	MinutesWatched int    `json:"minutes_watched"`
	Views          int    `json:"views"`
}

func (cfg *apiConfig) getChannelStats(w http.ResponseWriter, r *http.Request) {
	yt := YoutubeStats{
		Subscribers:    5,
		ChannelName:    "Bob's Stuff",
		MinutesWatched: 100,
		Views:          52,
	}

	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(cfg.ytApiKey))
	if err != nil {
		fmt.Println("failed to create service")
		panic(err)
	}

	call := yts.Channels.List([]string{"snippet, contentDetails, statistics"})
	call = call.ForUsername("Woohoojin")

	response, err := call.Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(response.Items[0].Snippet.Title)

	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(yt); err != nil {
		panic(err)
	}
}
