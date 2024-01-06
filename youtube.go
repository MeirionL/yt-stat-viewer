package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func getAccountYTChannel(tok string) (string, string, error) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tok},
	))

	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", "", err
	}

	call := yts.Channels.List([]string{"snippet"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		return "", "", err
	}

	return response.Items[0].Id, response.Items[0].Snippet.Title, nil
}

func (cfg *apiConfig) getYTChannelStats(w http.ResponseWriter, r *http.Request) {
	var response *youtube.ChannelListResponse
	isLive := "No"
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	channelString := chi.URLParam(r, "channel")
	lowerCaseChannel := strings.ToLower(channelString)

	user, err := cfg.DB.GetUserByChannelName(r.Context(), lowerCaseChannel)

	// handling channels not associated with authorised users
	if errors.Is(err, sql.ErrNoRows) {
		response, err = handleUnauthedChannel(ytApiKey, channelString)
		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("can't handle unauthed user %v: %v", channelString, err))
			return
		}

		if len(response.Items) > 0 {
			newChannel := Channel{
				Subscribers:    int(response.Items[0].Statistics.SubscriberCount),
				Title:          response.Items[0].Snippet.Title,
				Videos:         int(response.Items[0].Statistics.VideoCount),
				Views:          int(response.Items[0].Statistics.ViewCount),
				LastStreamTime: "",
				IsLive:         isLive,
				StreamTitle:    "",
			}
			cfg.handleChannelDeuplicates(newChannel)
			cfg.channels = append(cfg.channels, newChannel)
			RespondWithJSON(w, http.StatusOK, cfg.channels)
			return
		} else {
			RespondWithError(w, 400, "no items provided in response")
			return
		}
	}

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("can't get user by channel name %v: %v", channelString, err))
		return
	}

	response, broadcastResp, err := cfg.handleAuthedChannel(user)
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("can't handle authed user %v: %v", channelString, err))
		return
	}

	if len(broadcastResp.Items) > 0 {

		if broadcastResp.Items[0].Snippet.ActualStartTime != "" && broadcastResp.Items[0].Snippet.ActualEndTime == "" {
			isLive = "Yes!"
		}

		streamTime, err := handleTimeConversion(broadcastResp.Items[0].Snippet.ActualStartTime)
		if err != nil {
			RespondWithError(w, 400, "error converting time")
			return
		}

		newChannel := Channel{
			Subscribers:    int(response.Items[0].Statistics.SubscriberCount),
			Title:          response.Items[0].Snippet.Title,
			Videos:         int(response.Items[0].Statistics.VideoCount),
			Views:          int(response.Items[0].Statistics.ViewCount),
			LastStreamTime: streamTime,
			IsLive:         isLive,
			StreamTitle:    broadcastResp.Items[0].Snippet.Title,
		}
		cfg.handleChannelDeuplicates(newChannel)
		cfg.channels = append(cfg.channels, newChannel)
		RespondWithJSON(w, http.StatusOK, cfg.channels)
		return
	} else {
		newChannel := Channel{
			Subscribers:    int(response.Items[0].Statistics.SubscriberCount),
			Title:          response.Items[0].Snippet.Title,
			Videos:         int(response.Items[0].Statistics.VideoCount),
			Views:          int(response.Items[0].Statistics.ViewCount),
			LastStreamTime: "",
			IsLive:         isLive,
			StreamTitle:    broadcastResp.Items[0].Snippet.Title,
		}
		cfg.handleChannelDeuplicates(newChannel)
		cfg.channels = append(cfg.channels, newChannel)
		RespondWithJSON(w, http.StatusOK, cfg.channels)
		return
	}
}

func (cfg *apiConfig) handleAuthedChannel(u database.User) (*youtube.ChannelListResponse, *youtube.LiveBroadcastListResponse, error) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: u.AccessToken},
	))

	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, nil, err
	}

	call := yts.Channels.List([]string{"statistics, snippet"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == http.StatusUnauthorized {
			updatedUser, err := cfg.refreshAccessToken(u)
			if err != nil {
				return nil, nil, err
			}

			return cfg.handleAuthedChannel(updatedUser)
		} else {
			return nil, nil, err
		}
	}

	broadcastCall := yts.LiveBroadcasts.List([]string{"snippet"}).BroadcastType("all").MaxResults(10).Mine(true)
	broadcastResp, err := broadcastCall.Do()
	if err != nil {
		return nil, nil, err
	}

	return response, broadcastResp, nil
}

func handleUnauthedChannel(ytApiKey, channelString string) (*youtube.ChannelListResponse, error) {
	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithAPIKey(ytApiKey))
	if err != nil {
		return nil, err
	}

	firstCall := yts.Search.List([]string{"snippet"}).Type("channel").Q(channelString).MaxResults(1)
	firstResponse, err := firstCall.Do()
	if err != nil {
		return nil, err
	}

	if len(firstResponse.Items) == 0 {
		return nil, fmt.Errorf("no channel returned for %v", channelString)
	}

	secondCall := yts.Channels.List([]string{"statistics, snippet"}).Id(firstResponse.Items[0].Id.ChannelId)
	secondResponse, err := secondCall.Do()
	if err != nil {
		return nil, err
	}

	return secondResponse, nil
}

func handleTimeConversion(timeString string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	if err != nil {
		return "", err
	}

	formattedTime := parsedTime.Format("02/01/2006 15:04")

	return formattedTime, nil
}

func (cfg *apiConfig) handleChannelDeuplicates(c Channel) {
	for i, channel := range cfg.channels {
		if channel.Title == c.Title {
			cfg.channels = append(cfg.channels[:i], cfg.channels[i+1:]...)
			return
		}
	}
}
