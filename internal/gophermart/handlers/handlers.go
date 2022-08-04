package handlers

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(auth services.Authenticator) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	h := NewHandler(auth)

	r.Post("/api/user/register", h.Register)

	return r
}

type Handler struct {
	Mux  *chi.Mux
	auth services.Authenticator
}

func NewHandler(auth services.Authenticator) *Handler {
	return &Handler{
		Mux:  chi.NewMux(),
		auth: auth,
	}
}

func getDecompressedReader(r *http.Request) (io.Reader, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}
