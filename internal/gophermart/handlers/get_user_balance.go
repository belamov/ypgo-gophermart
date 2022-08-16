package handlers

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	userWithdrawalsAmount, err := h.balanceProcessor.GetUserTotalWithdrawAmount(userID)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user balance handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userAccrualAmount, err := h.balanceProcessor.GetUserTotalAccrualAmount(userID)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user balance handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
	balance := userAccrualAmount - userWithdrawalsAmount
	result := response{
		Current:   math.Round(balance*100) / 100, //nolint:gomnd
		Withdrawn: userWithdrawalsAmount,
	}

	out, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in get user balance handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(out); err != nil {
		log.Error().Err(err).Msg("unexpected error in get user balance handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
