package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

func (h *Handler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	userWithdrawals, err := h.balanceProcessor.GetUserWithdrawals(userID)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user withdraws handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userWithdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type WithdrawalResponse struct {
		Order     string  `json:"order"`
		Sum       float64 `json:"sum"`
		CreatedAt string  `json:"processed_at"`
	}
	result := make([]WithdrawalResponse, len(userWithdrawals))

	for i, withdrawal := range userWithdrawals {
		result[i] = WithdrawalResponse{
			Order:     strconv.Itoa(withdrawal.OrderID),
			Sum:       withdrawal.WithdrawalAmount,
			CreatedAt: withdrawal.CreatedAt.Format(time.RFC3339),
		}
	}

	out, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user withdraws handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(out); err != nil {
		log.Error().Err(err).Msg("unexpected error in get user withdraws handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
