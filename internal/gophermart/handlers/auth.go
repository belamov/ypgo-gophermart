package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var credentials models.Credentials

	if err := json.NewDecoder(reader).Decode(&credentials); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if err := credentials.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid credentials: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.auth.Register(credentials)
	if err != nil {
		http.Error(w, "cannot register user", http.StatusInternalServerError)
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	addAuthHeader(token, w)
}

func addAuthHeader(token string, w http.ResponseWriter) {
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
}
