package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func (cfg *apiConfig) getStats(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, cfg.channels)
}

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

	sort.Slice(response.Items, func(i, j int) bool {
		return response.Items[i].Statistics.SubscriberCount > response.Items[j].Statistics.SubscriberCount
	})

	fmt.Printf("\nname of most subscribed YT channel: %v\n", response.Items[0].Snippet.Title)
	fmt.Printf("\nID of most subscribed YT channel: %v\n", response.Items[0].Id)

	return response.Items[0].Id, response.Items[0].Snippet.Title, nil
}

func (cfg *apiConfig) getYTChannelStats(w http.ResponseWriter, r *http.Request) {
	var user database.User
	var response *youtube.ChannelListResponse
	var err error
	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	// I hate having these here
	channelString := chi.URLParam(r, "channel")

	user, err = cfg.DB.GetUserByChannelName(r.Context(), channelString)

	// handling channels not associated with authorised users
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("\nunatuthenticated user!!!\n")
		response, err = handleUnauthedChannel(ytApiKey, channelString)
		if err != nil {
			fmt.Println(err)
			RespondWithError(w, 400, fmt.Sprintf("can't handle unauthed user %v: %v", channelString, err))
			return
		}

		fmt.Printf("\nchannel name: %v\n", response.Items[0].Snippet.Title)
		fmt.Printf("\nchannel subscriber count: %v\n", response.Items[0].Statistics.SubscriberCount)
		fmt.Printf("\nchannel views count: %v\n", response.Items[0].Statistics.ViewCount)
		fmt.Printf("\nchannels liked playlists: %v\n", response.Items[0].ContentDetails.RelatedPlaylists.Likes)
		fmt.Printf("\nchannel video count: %v\n", response.Items[0].Statistics.VideoCount)

		if len(response.Items) > 0 {
			newChannel := Channel{
				Subscribers: int(response.Items[0].Statistics.SubscriberCount),
				Title:       response.Items[0].Snippet.Title,
				Views:       int(response.Items[0].Statistics.ViewCount),
				Platform:    "YouTube",
			}
			cfg.channels = append(cfg.channels, newChannel)
			RespondWithJSON(w, http.StatusOK, cfg.channels)
			return
		} else {
			RespondWithError(w, 400, "no items provided in response")
			return
		}
	} else if err != nil {
		fmt.Println(err)
		RespondWithError(w, 400, fmt.Sprintf("can't get user by channel name %v: %v", channelString, err))
		return
	}

	fmt.Printf("\nauthenticated user!!!\n")
	response, err = cfg.handleAuthedChannel(user)
	if err != nil {
		fmt.Println(err)
		RespondWithError(w, 400, fmt.Sprintf("can't handle authed user %v: %v", channelString, err))
		return
	}

	fmt.Printf("\nchannel name: %v\n", response.Items[0].Snippet.Title)
	fmt.Printf("\nchannel subscriber count: %v\n", response.Items[0].Statistics.SubscriberCount)
	fmt.Printf("\nchannel views count: %v\n", response.Items[0].Statistics.ViewCount)
	fmt.Printf("\nchannels liked playlists: %v\n", response.Items[0].ContentDetails.RelatedPlaylists.Likes)
	fmt.Printf("\nchannel video count: %v\n", response.Items[0].Statistics.VideoCount)

}

func (cfg *apiConfig) handleAuthedChannel(u database.User) (*youtube.ChannelListResponse, error) {
	httpClient := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: u.AccessToken},
	))

	ctx := context.Background()
	yts, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	fmt.Println("starting normal call")
	call := yts.Channels.List([]string{"contentDetails, statistics, snippet"}).Mine(true)
	response, err := call.Do()
	if err != nil {
		if apiErr, ok := err.(*googleapi.Error); ok && apiErr.Code == http.StatusUnauthorized {
			// code that tries to refresh access token
			updatedUser, err := cfg.refreshAccessToken(u)
			if err != nil {
				fmt.Println("maybe expired refresh token")
				return nil, err
			}

			return cfg.handleAuthedChannel(updatedUser)
		} else {
			// Handle other non-nil errors as usual
			fmt.Printf("\nerror is mooboo: %v\n", err)
			return nil, err
		}

	}

	fmt.Printf("\nerror is: %v\n", err)

	// fmt.Println("starting members call")
	// memberCall := yts.Members.List([]string{"snippet"})
	// memberResp, err := memberCall.Do()
	// if err != nil {
	// 	return nil, err
	// }

	// if memberResp.PageInfo != nil {
	// 	fmt.Printf("\ntotal members for channel: %v\n", memberResp.PageInfo.TotalResults)
	// } else {
	// 	fmt.Print("\nno channel members were found\n")
	// }

	// PERMISSIONS TO STREAM ARE AT 10:40
	// think i have to give up on members
	fmt.Println("starting broadcast call")
	// broadcastCall := yts.LiveBroadcasts.List([]string{"snippet, statistics"}).BroadcastType("all").MaxResults(10).Mine(true)
	// broadcastResp, err := broadcastCall.Do()
	// if err != nil {
	// 	return nil, err
	// }

	// livestreams := broadcastResp.Items
	// fmt.Printf("\nhere are the channels 10 most recent stream stats:\n")
	// if len(livestreams) > 0 {
	// 	for _, livestream := range livestreams {
	// 		viewers := livestream.Statistics.ConcurrentViewers
	// 		title := livestream.Snippet.Title
	// 		fmt.Printf("Livestream: %s | Viewers: %v\n", title, viewers)
	// 	}
	// }

	return response, nil
}

func (cfg *apiConfig) refreshAccessToken(u database.User) (database.User, error) {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	ctx := context.Background()

	config := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost:8080/auth/callback",
	}

	token := &oauth2.Token{
		RefreshToken: u.RefreshToken,
	}

	newToken, err := config.TokenSource(ctx, token).Token()
	if err != nil {
		fmt.Printf("\ncouldn't get access token using refresh token: %v\n", err)
		return database.User{}, err
	}

	if newToken.AccessToken != u.AccessToken && newToken.RefreshToken == u.RefreshToken {
		newUser, err := cfg.DB.UpdateAccessToken(ctx, database.UpdateAccessTokenParams{
			ID:          u.ID,
			UpdatedAt:   time.Now().UTC(),
			AccessToken: newToken.AccessToken,
		})
		if err != nil {
			fmt.Printf("\ncouldn't update new access token to user: %v\n", err)
			return database.User{}, err
		}

		return newUser, nil
	} else {
		fmt.Println("access token isn't new/refresh token has changed. Can't confirm authorization")
		return database.User{}, err
	}
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

	secondCall := yts.Channels.List([]string{"contentDetails, statistics, snippet"}).Id(firstResponse.Items[0].Id.ChannelId)

	secondResponse, err := secondCall.Do()
	if err != nil {
		return nil, err
	}

	return secondResponse, nil
}
