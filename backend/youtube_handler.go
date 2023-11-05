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
	Subscribers int    `json:"subscribers"`
	ChannelName string `json:"channel_name"`
	Views       int    `json:"views"`
}

func (cfg *apiConfig) getChannelStats(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(cfg.ytApiKey))
	if err != nil {
		fmt.Println("failed to create service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	firstCall := yts.Search.List([]string{"snippet"}).Type("channel").Q("Mogul Mail")

	firstResponse, err := firstCall.Do()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelID := firstResponse.Items[0].Id.ChannelId

	secondCall := yts.Channels.List([]string{"snippet, contentDetails, statistics"})
	secondCall = secondCall.Id(channelID)

	secondResponse, err := secondCall.Do()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yt := YoutubeStats{}
	if len(firstResponse.Items) > 0 && len(secondResponse.Items) > 0 {
		val := secondResponse.Items[0]
		yt = YoutubeStats{
			Subscribers: int(val.Statistics.SubscriberCount),
			ChannelName: val.Snippet.Title,
			Views:       int(val.Statistics.ViewCount),
		}
	}

	// Broken code that should work for all channels
	//
	// call := yts.Channels.List([]string{"snippet, contentDetails, statistics"})
	// call = call.ForUsername("Pewdiepie")
	// response, err := call.Do()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
	// 	"and it has %d views.",
	// 	response.Items[0].Id,
	// 	response.Items[0].Snippet.Title,
	// 	response.Items[0].Statistics.ViewCount))

	w.Header().Set("Conten-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(yt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
