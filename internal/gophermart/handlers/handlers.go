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

func NewRouter(auth services.Auth, orders services.OrdersProcessorInterface) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	h := NewHandler(auth, orders)
	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware())
		r.Post("/api/user/orders", h.AddOrder)
	})

	return r
}

type Handler struct {
	Mux    *chi.Mux
	auth   services.Auth
	orders services.OrdersProcessorInterface
}

func NewHandler(auth services.Auth, orders services.OrdersProcessorInterface) *Handler {
	return &Handler{
		Mux:    chi.NewMux(),
		auth:   auth,
		orders: orders,
	}
}

func getDecompressedReader(r *http.Request) (io.Reader, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}
