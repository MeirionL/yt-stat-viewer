package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type Channel struct {
	Subscribers int    `json:"subscribers"`
	Title       string `json:"title"`
	Views       int    `json:"views"`
	Platform    string `json:"platform"`
}

type apiConfig struct {
	ytApiKey   string
	channels   []Channel
	query      *string
	maxResults *int64
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the enviroment")
	}

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	if ytApiKey == "" {
		log.Fatal("youtube api key is not found in the enviroment")
	}

	cfg := apiConfig{
		ytApiKey: ytApiKey,
		channels: []Channel{},
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/healthz", HandlerReadiness)
	router.Get("/err", HandlerErr)
	router.Get("/youtube/stats", cfg.getStats)
	router.Get("/youtube/stats/{channel}", cfg.getYTChannelStats)

	fs := http.FileServer(http.Dir("."))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	log.Fatal(srv.ListenAndServe())
}
