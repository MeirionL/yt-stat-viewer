package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/nicklaw5/helix/v2"
)

func (cfg *apiConfig) getTwitchChannelStats(w http.ResponseWriter, r *http.Request) {
	channelString := chi.URLParam(r, "channel")

	client, err := helix.NewClient(&helix.Options{
		ClientID:      cfg.twitchClientID,
		RateLimitFunc: rateLimitCallback,
	})
	if err != nil {
		fmt.Println("failed to create service")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp1, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{channelString},
	})
	if err != nil {
		fmt.Printf("failed to get users: %s\n", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("\n%vs display name: %v\n", channelString, resp1.Data.Users[0].DisplayName)
	fmt.Printf("\n%vs display name: %v\n", channelString, resp1.Data.Users[0].ID)
	fmt.Printf("\n%vs display name: %v\n", channelString, resp1.Data.Users[0].ViewCount)

	// resp1, err := client.GetChannelInformation(&helix.GetChannelInformationParams{
	// 	BroadcasterIDs: []string{"71092938"},
	// })
	// if err != nil {
	// 	fmt.Printf("failed to get users: %s\n", err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Println(resp1.Data.Channels[0].BroadcasterID)
	// fmt.Println(resp1.Data.Channels[0].Title)
	// fmt.Println(resp1.Data.Channels[0].BroadcasterLanguage)
	// fmt.Printf("my channel response: %v\n\n\n", resp1)

	fmt.Printf("Status code: %d\n", resp1.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp1.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp1.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp1.GetRateLimitReset())
	// fmt.Println(resp1.Data.Users[0].ID)
	// fmt.Println(resp1.Data.Users[0].CreatedAt.Time)
	// fmt.Println(resp1.Data.Users[0].Email)
	// fmt.Println(resp1.Data.Users[0])

	// if len(resp1.Data.Channels) == 0 {
	// 	fmt.Printf("No users returned for the given channel %v\n", channelString)
	// 	return
	// }

	// for _, user := range resp1.Data.Channels {
	// 	fmt.Printf("\n%v is the current channel\n", user)
	// }

	// resp2, err := client.GetSubscriptions(&helix.SubscriptionsParams{
	// 	BroadcasterID: resp1.Data.Channels[0].BroadcasterID,
	// })
	// if err != nil {
	// 	fmt.Println("failed to get subscribers")
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// fmt.Printf("Total Twitch subscriptions to %v is: %v\n", channelString, resp2.Data.Subscriptions)

	// newChannel := Channel{
	// 	Subscribers: resp2.Data.Total,
	// 	Title:       resp1.Data.Channels[0].Title,
	// 	Views:       5000,
	// 	Platform:    "Twitch",
	// }
	// cfg.channels = append(cfg.channels, newChannel)

	// RespondWithJSON(w, http.StatusOK, cfg.channels)
}

func (cfg *apiConfig) getAccountTwitchChannel(u goth.User) (string, error) {
	client, err := helix.NewClient(&helix.Options{
		ClientID:        cfg.twitchClientID,
		RateLimitFunc:   rateLimitCallback,
		UserAccessToken: u.AccessToken,
	})
	if err != nil {
		fmt.Println("failed to create service")
		return "", err
	}

	resp, err := client.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{u.UserID},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Data.Channels) == 0 {
		fmt.Printf("\nAin't got no channels: \n%v", u.UserID)
		return "", err
	}

	fmt.Printf("\ntwitch channel found for user: %v\n", resp.Data.Channels[0].BroadcasterName)
	fmt.Printf("\nID of twitch channel: %v\n", resp.Data.Channels[0].BroadcasterID)
	fmt.Printf("\nID of user: %v\n", u.UserID)

	return resp.Data.Channels[0].BroadcasterName, nil
}

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	reset64 := int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			fmt.Printf("Waiting on rate limit to pass before sending next request (%d seconds)\n", timeDiff)
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}
