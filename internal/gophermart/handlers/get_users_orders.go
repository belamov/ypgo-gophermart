package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func (h *Handler) GetUsersOrders(w http.ResponseWriter, r *http.Request) {
	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	usersOrders, err := h.ordersManager.GetUsersOrders(userID)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user orders handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(usersOrders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type OrderResponse struct {
		Number     string  `json:"number"`
		Status     string  `json:"status"`
		Accrual    float64 `json:"accrual,omitempty"`
		UploadedAt string  `json:"uploaded_at"`
	}

	result := make([]OrderResponse, len(usersOrders))

	for i, order := range usersOrders {
		result[i] = OrderResponse{
			Number:     strconv.Itoa(order.ID),
			Status:     order.Status.String(),
			Accrual:    math.Round(order.Accrual*100) / 100, //nolint:gomnd
			UploadedAt: order.UploadedAt.Format(time.RFC3339),
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user orders handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(out); err != nil {
		log.Error().Err(err).Msg("unexpected error in get user orders handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
