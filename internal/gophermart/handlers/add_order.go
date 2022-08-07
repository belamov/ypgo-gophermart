package handlers

import (
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) AddOrder(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rawOrderID, err := io.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orderID, err := strconv.Atoi(string(rawOrderID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err = h.orders.ValidateOrderId(orderID); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	userID := h.auth.GetUserId(r)
	if userID == 0 {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	err = h.orders.AddOrder(orderID, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
