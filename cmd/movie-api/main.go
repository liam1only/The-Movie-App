package main

import (
	"log"
	"net/http"
	"os"

	"example.com/movie-api/internal/httpapi"
	"example.com/movie-api/internal/tmdb"
)

func main() {
	moviesClient := tmdb.NewClient(os.Getenv("TMDB_API_KEY"))
	handler := httpapi.NewHandler(moviesClient)

	log.Println("movie API listening on :8080")
	if err := http.ListenAndServe(":8080", handler.Routes()); err != nil {
		log.Fatal(err)
	}
}
