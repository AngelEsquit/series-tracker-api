package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Rutas de la API
func setupRoutes(router *mux.Router) {
	router.HandleFunc("/api/series", createSeries).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/series", getAllSeries).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/series/{id}", getSeriesByID).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/series/{id}", updateSeries).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/series/{id}", deleteSeries).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/series/{id}/status", updateStatus).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/series/{id}/episode", updateEpisode).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/series/{id}/upvote", upVote).Methods("PATCH", "OPTIONS")
	router.HandleFunc("/api/series/{id}/downvote", downVote).Methods("PATCH", "OPTIONS")
}

// enableCORS habilita CORS para todas las rutas de la API
func enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
