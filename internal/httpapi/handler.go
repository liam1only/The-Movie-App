package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/movie-api/internal/tmdb"
)

type MovieSearcher interface {
	SearchMovies(ctx context.Context, query string) (tmdb.SearchResponse, error)
}

type Handler struct {
	movies MovieSearcher
}

func NewHandler(movies MovieSearcher) *Handler {
	return &Handler{movies: movies}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /movies/search", h.searchMovies)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) searchMovies(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		writeError(w, http.StatusBadRequest, "query parameter q is required")
		return
	}

	result, err := h.movies.SearchMovies(r.Context(), query)
	if errors.Is(err, tmdb.ErrMissingAPIKey) {
		writeError(w, http.StatusServiceUnavailable, "TMDB is not configured")
		return
	}
	if err != nil {
		writeError(w, http.StatusBadGateway, "movie search is temporarily unavailable")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
