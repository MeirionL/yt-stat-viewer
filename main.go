package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	ytApiKey string
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

	fsHandler := cfg.handlerHTML(http.StripPrefix("/app", http.FileServer(http.Dir("client/react_content"))))
	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	router.Get("/healthz", HandlerReadiness)
	router.Get("/err", HandlerErr)
	router.Get("/youtube/channel/stats", cfg.getChannelStats)

	fs := http.FileServer(http.Dir("client/react_content"))
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

func (cfg *apiConfig) handlerHTML(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}
