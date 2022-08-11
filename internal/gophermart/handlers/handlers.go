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

func NewRouter(
	auth services.Auth,
	ordersProcessor services.OrdersProcessorInterface,
	balanceProcessor services.BalanceProcessorInterface,
) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(flate.BestSpeed))

	h := NewHandler(auth, ordersProcessor, balanceProcessor)
	r.Post("/api/user/register", h.Register)
	r.Post("/api/user/login", h.Login)
	r.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware())
		r.Post("/api/user/orders", h.AddOrder)
		r.Get("/api/user/orders", h.GetUsersOrders)
		r.Post("/api/user/balance/withdraw", h.RegisterWithdraw)
		r.Get("/api/user/withdrawals", h.GetUserWithdrawals)
	})

	return r
}

type Handler struct {
	Mux              *chi.Mux
	auth             services.Auth
	ordersProcessor  services.OrdersProcessorInterface
	balanceProcessor services.BalanceProcessorInterface
}

func NewHandler(
	auth services.Auth,
	ordersProcessor services.OrdersProcessorInterface,
	balanceProcessor services.BalanceProcessorInterface,
) *Handler {
	return &Handler{
		Mux:              chi.NewMux(),
		auth:             auth,
		ordersProcessor:  ordersProcessor,
		balanceProcessor: balanceProcessor,
	}
}

func getDecompressedReader(r *http.Request) (io.Reader, error) {
	if r.Header.Get("Content-Encoding") == "gzip" {
		return gzip.NewReader(r.Body)
	}
	return r.Body, nil
}
