package handlers

import (
	"compress/flate"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	// h := NewHandler()
	//
	// r.Get("/ping", h.Ping)

	return r
}

type Handler struct {
	Mux *chi.Mux
}

// func NewHandler() *Handler {
//	return &Handler{
//		Mux: chi.NewMux(),
//	}
//}
//
// func getDecompressedReader(r *http.Request) (io.Reader, error) {
//	if r.Header.Get("Content-Encoding") == "gzip" {
//		return gzip.NewReader(r.Body)
//	}
//	return r.Body, nil
//}
