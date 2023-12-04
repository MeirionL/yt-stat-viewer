package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

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
	DB                 *database.Queries
	ytApiKey           string
	googleClientID     string
	googleClientSecret string
	twitchClientID     string
	channels           []Channel
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

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	if googleClientID == "" {
		log.Fatal("google client id is not found in enviroment")
	}

	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if googleClientSecret == "" {
		log.Fatal("google client secret is not found in enviroment")
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
		DB:                 database.New(conn),
		ytApiKey:           ytApiKey,
		googleClientID:     googleClientID,
		googleClientSecret: googleClientSecret,
		twitchClientID:     twitchClientID,
		channels:           []Channel{},
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

	router.Post("/users", cfg.handlerCreateUser)

	router.Get("/stats", cfg.getStats)
	router.Get("/stats/YouTube/{channel}", cfg.getYTChannelStats)
	router.Get("/stats/Twitch/{channel}", cfg.getTwitchChannelStats)
	router.Get("/stats/Kick/{channel}", cfg.getKickChannelStats)

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
