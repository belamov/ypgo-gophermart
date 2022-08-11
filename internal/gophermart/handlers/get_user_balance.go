package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID := h.auth.GetUserID(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	userWithdrawalsAmount, err := h.balanceProcessor.GetUserTotalWithdrawAmount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userAccrualAmount, err := h.balanceProcessor.GetUserTotalAccrualAmount(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type response struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
	result := response{
		Current:   userAccrualAmount - userWithdrawalsAmount,
		Withdrawn: userWithdrawalsAmount,
	}

	out, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(out); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
