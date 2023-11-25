package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (cfg *apiConfig) getStats(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, cfg.channels)
}

func (cfg *apiConfig) getYTChannelStats(w http.ResponseWriter, r *http.Request) {
	channelString := chi.URLParam(r, "channel")
	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(cfg.ytApiKey))
	if err != nil {
		fmt.Println("failed to create service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	firstCall := yts.Search.List([]string{"snippet"}).Type("channel").Q(channelString).MaxResults(5)

	firstResponse, err := firstCall.Do()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channelID := firstResponse.Items[0].Id.ChannelId

	secondCall := yts.Channels.List([]string{"contentDetails, statistics"})
	secondCall = secondCall.Id(channelID)

	secondResponse, err := secondCall.Do()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(firstResponse.Items) > 0 && len(secondResponse.Items) > 0 {
		newChannel := Channel{
			Subscribers: int(secondResponse.Items[0].Statistics.SubscriberCount),
			Title:       firstResponse.Items[0].Snippet.ChannelTitle,
			Views:       int(secondResponse.Items[0].Statistics.ViewCount),
			Platform:    "YouTube",
		}
		cfg.channels = append(cfg.channels, newChannel)
	}

	// Code that only works for legacy channels
	//
	// call := yts.Search.List([]string{"snippet, contentDetails, statistics"})
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

	RespondWithJSON(w, http.StatusOK, cfg.channels)
}
