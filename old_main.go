package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
	Body  string `json:"body"`
}

func old_main() {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	todos := []Todo{}

	router.Get("/err", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		dat, err := json.Marshal(map[string]string{"status": "ok"})
		if err != nil {
			log.Printf("error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
		w.Write(dat)
	})

	router.Post("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		todo := &Todo{}

		// Decode JSON from request body into the Todo struct
		err := json.NewDecoder(r.Body).Decode(&todo)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("couldn't decode JSON: %v", err))
			return
		}

		// Assign an ID to the todo
		todo.ID = len(todos) + 1

		// Append the todo to the todos slice
		todos = append(todos, *todo)

		// Respond with the updated list of todos including the new todo
		RespondWithJSON(w, http.StatusOK, todos)
	})

	router.Patch("/api/todos/{id}/done", func(w http.ResponseWriter, r *http.Request) {
		idString := chi.URLParam(r, "id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("couldn't convert todo id to int: %v", err))
			return
		}

		for i, t := range todos {
			if t.ID == id {
				todos[i].Done = true
				break
			}
		}

		RespondWithJSON(w, http.StatusOK, todos)
	})

	router.Get("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		RespondWithJSON(w, http.StatusOK, todos)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + "8080",
	}

	log.Printf("Server starting on port 8080")
	log.Fatal(srv.ListenAndServe())
}
