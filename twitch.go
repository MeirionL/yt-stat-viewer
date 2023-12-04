package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

	fmt.Printf("%v\n", resp1.Header)
	fmt.Printf("Status code: %d\n", resp1.StatusCode)
	fmt.Printf("Rate limit: %d\n", resp1.GetRateLimit())
	fmt.Printf("Rate limit remaining: %d\n", resp1.GetRateLimitRemaining())
	fmt.Printf("Rate limit reset: %d\n\n", resp1.GetRateLimitReset())

	if len(resp1.Data.Users) == 0 {
		fmt.Println("No users returned for the given channel")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, user := range resp1.Data.Users {
		fmt.Printf("%v is the current channel", user)
	}

	resp2, err := client.GetSubscriptions(&helix.SubscriptionsParams{
		BroadcasterID: resp1.Data.Users[0].ID,
	})
	if err != nil {
		fmt.Println("failed to get subscribers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("Total Twitch subscriptions to %v is: %v", channelString, resp2.Data.Total)

	newChannel := Channel{
		Subscribers: 500,
		Title:       channelString,
		Views:       5000,
		Platform:    "Twitch",
	}
	cfg.channels = append(cfg.channels, newChannel)

	RespondWithJSON(w, http.StatusOK, cfg.channels)
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
