// internal/api/router.go
package api

import (
    "net/http"

    "github.com/go-chi/chi/v5" // or "github.com/gorilla/mux"
    // If you don't have chi installed yet, run: go get github.com/go-chi/chi/v5
)

func NewRouter() http.Handler {
    r := chi.NewRouter()
    // Add endpoints
    r.Get("/events", GetEventsHandler)
    return r
}
