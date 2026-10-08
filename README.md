# Movie API

A small Go API for finding movies and, as the project grows, tracking upcoming releases, watchlists, and reminders. Movie information is provided by [The Movie Database (TMDB)](https://www.themoviedb.org/).

## Current features

- `GET /health` — reports whether the API server is running.
- `GET /movies/search?q=title` — searches TMDB for movies matching a title.

Search results include TMDB's pagination details and movie fields such as the TMDB ID, title, overview, poster path, and release date.

## Requirements

- Go (use the version specified in `go.mod`)
- A TMDB API key

Create a TMDB account and request an API key in your account's API settings. Keep the key private; the server reads it from the `TMDB_API_KEY` environment variable.

## Run locally

In the project root, set the key and start the server:

```sh
export TMDB_API_KEY="your_tmdb_api_key"
go run ./cmd/movie-api
```

The server listens on port `8080`.

## Try the API

Check that the server is running:

```sh
curl http://localhost:8080/health
```

Search for a movie:

```sh
curl "http://localhost:8080/movies/search?q=Inception"
```

The search endpoint returns `400 Bad Request` when `q` is missing, `503 Service Unavailable` when `TMDB_API_KEY` is not configured, and `502 Bad Gateway` if TMDB cannot complete the request.

## Project layout

```text
cmd/movie-api/       API server entry point
internal/httpapi/    HTTP routes and handlers
internal/tmdb/       TMDB API client
```
