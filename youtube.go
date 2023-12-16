package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	googleClientID     = os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	ytApiKey           = os.Getenv("YOUTUBE_API_KEY")
)

func (cfg *apiConfig) getStats(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, cfg.channels)
}

func (cfg *apiConfig) handleYoutubeAuth(w http.ResponseWriter, r *http.Request) {
	// conf := &oauth2.Config{
	// 	ClientID:     googleClientID,
	// 	ClientSecret: googleClientSecret,
	// 	RedirectURL:  "http://localhost:5173",
	// 	Scopes: []string{"https://www.googleapis.com/auth/youtube.readonly",
	// 		"https://www.googleapis.com/auth/youtube.channel-memberships.creator",
	// 		"https://www.googleapis.com/auth/yt-analytics.readonly"},
	// 	Endpoint: google.Endpoint,
	// }

	// if code := r.URL.Query().Get("code"); code != "" {
	// 	// Exchange authorization code for an access token
	// 	token, err := conf.Exchange(context.Background(), code)
	// 	if err != nil {
	// 		// Handle error
	// 		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// Access token can now be used to make authorized requests to the YouTube API
	// 	// Store the access token securely for future API requests

	// 	// Example: Print the access token
	// 	fmt.Printf("Access Token: %s", token.AccessToken)
	// 	return
	// }

	// url := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	// http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func getAccountYTChannel(u goth.User) (string, error) {
	tokenString := u.AccessToken

	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tokenString},
	))

	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", err
	}

	call := yts.Channels.List([]string{"snippet"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	var channelNames []string
	for _, channel := range response.Items {
		channelNames = append(channelNames, channel.Snippet.Title)
	}

	fmt.Printf("\nall YT channels found for user: %v\n", channelNames)
	fmt.Printf("\nID of first YT channel: %v\n", response.Items[0].Id)
	fmt.Printf("\nID of user: %v\n", u.IDToken)

	return channelNames[0], nil
}

func (cfg *apiConfig) getYTChannelStats(w http.ResponseWriter, r *http.Request) {
	channelString := chi.URLParam(r, "channel")
	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(ytApiKey))
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
