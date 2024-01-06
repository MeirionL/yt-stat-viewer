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
	Title          string `json:"title"`
	Subscribers    int    `json:"subscribers"`
	Videos         int    `json:"videos"`
	Views          int    `json:"views"`
	LastStreamTime string `json:"last_stream_time"`
	IsLive         string `json:"is_live"`
	StreamTitle    string `json:"stream_title"`
}

type apiConfig struct {
	DB       *database.Queries
	ytApiKey string
	channels []Channel
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

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("unable to establish database connection", err)
	}

	cfg := apiConfig{
		DB:       database.New(conn),
		ytApiKey: ytApiKey,
		channels: []Channel{},
	}

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

	router.Get("/stats/{channel}", cfg.getYTChannelStats)

	router.Get("/auth/callback", cfg.get_google_token)
	router.Get("/auth", google_auth)

	router.Get("/logout/{id}", cfg.logout)

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
