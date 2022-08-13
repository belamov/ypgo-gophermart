package handlers

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(orderManager services.OrderManagementInterface) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	h := NewHandler(orderManager)

	r.Post("/api/orders", h.RegisterOrder)

	return r
}

type Handler struct {
	Mux          *chi.Mux
	orderManager services.OrderManagementInterface
}

func NewHandler(orderManager services.OrderManagementInterface) *Handler {
	return &Handler{
		Mux:          chi.NewMux(),
		orderManager: orderManager,
	}
}

func getDecompressedReader(r *http.Request) (io.Reader, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}
