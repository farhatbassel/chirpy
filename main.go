package main

import (
	"fmt"
	"net/http"

	db "farhatbassel/chirpy/internal/database"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()

	data, err := db.NewDB("database.json")

	if err != nil {
		panic(err)
	}

	apiConfig := apiConfig{
		db: *data,
	}

	// Static Routes
	router.Handle("/app/",
		apiConfig.middlewareMetricsInc(
			http.StripPrefix("/app/", http.FileServer(http.Dir("."))),
		),
	)
	router.Handle("/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("/assets"))),
	)

	// Admin routes
	router.Get("/admin/metrics", apiConfig.displayNumberOfHits)

	// API Routes
	router.HandleFunc("/api/reset", apiConfig.resetNumberOfHits)
	router.Get("/api/healthz", checkForHealth)
	router.Post("/api/chirps", apiConfig.createChirp)
	router.Get("/api/chirps", apiConfig.getChirps)

	corsMux := middlewareCors(router)

	fmt.Println("Starting server on port 8080...")
	server := http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	server.ListenAndServe()
}

func checkForHealth(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}
