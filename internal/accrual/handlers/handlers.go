package handlers

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/accrual/services"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(orderManager services.OrderManagementInterface, rewardsStorage storage.RewardsStorage) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	h := NewHandler(orderManager, rewardsStorage)

	r.Post("/api/orders", h.RegisterOrder)
	r.Post("/api/goods", h.RegisterReward)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Throttle(100)) //nolint:gomnd
		r.Get("/api/orders/{order}", h.GetOrderInfo)
	})

	return r
}

type Handler struct {
	Mux            *chi.Mux
	orderManager   services.OrderManagementInterface
	rewardsStorage storage.RewardsStorage
}

func NewHandler(orderManager services.OrderManagementInterface, rewardsStorage storage.RewardsStorage) *Handler {
	return &Handler{
		Mux:            chi.NewMux(),
		orderManager:   orderManager,
		rewardsStorage: rewardsStorage,
	}
}

func getDecompressedReader(r *http.Request) (io.Reader, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}
