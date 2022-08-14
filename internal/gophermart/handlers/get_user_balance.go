package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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
		log.Println("unexpected error in get user balance handler:")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userAccrualAmount, err := h.balanceProcessor.GetUserTotalAccrualAmount(userID)
	if err != nil {
		log.Println("unexpected error in get user balance handler:")
		log.Println(err.Error())
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
		log.Println("unexpected error in get user balance handler:")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write(out); err != nil {
		log.Println("unexpected error in get user balance handler:")
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
