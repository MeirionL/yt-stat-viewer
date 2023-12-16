package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/MeirionL/stream_stats/backend/internal/auth"
	"github.com/MeirionL/stream_stats/backend/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Channel struct {
	Subscribers int    `json:"subscribers"`
	Title       string `json:"title"`
	Views       int    `json:"views"`
	Platform    string `json:"platform"`
}

type apiConfig struct {
	DB             *database.Queries
	ytApiKey       string
	twitchClientID string
	channels       []Channel
}

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("port is not found in the enviroment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("database url is not found in the enviroment")
	}

	ytApiKey := os.Getenv("YOUTUBE_API_KEY")
	if ytApiKey == "" {
		log.Fatal("youtube api key is not found in the enviroment")
	}

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		log.Fatal("twitch client id is not found in enviroment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("unable to establish database connection", err)
	}

	cfg := apiConfig{
		DB:             database.New(conn),
		ytApiKey:       ytApiKey,
		twitchClientID: twitchClientID,
		channels:       []Channel{},
	}

	auth.NewAuth()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/healthz", HandlerReadiness)
	router.Get("/err", HandlerErr)

	router.Post("/users", cfg.handlerCreateUser)
	router.Get("/users", cfg.handlerGetUsers)
	// router.Put("/users", cfg.handlerUpdateUser)
	// router.Delete("/users", cfg.handlerDeleteUser)

	router.Get("/stats", cfg.getStats)
	// router.Get("/stats/YouTube/auth", cfg.handleYoutubeAuth)
	router.Get("/stats/YouTube/{channel}", cfg.getYTChannelStats)
	router.Get("/stats/Twitch/{channel}", cfg.getTwitchChannelStats)

	router.Get("/auth/{provider}", cfg.handlerAuthLogin)
	router.Get("/logout/{provider}", cfg.handlerAuthLogout)

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
