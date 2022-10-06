package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/gophermart/models"
	"github.com/belamov/ypgo-gophermart/internal/gophermart/services"
	"github.com/rs/zerolog/log"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var credentials models.Credentials

	if err = json.NewDecoder(reader).Decode(&credentials); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if err = credentials.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid credentials: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.auth.Register(credentials)

	var loginTakenError *services.LoginTakenError
	if errors.As(err, &loginTakenError) {
		http.Error(w, "login is taken", http.StatusConflict)
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register handler")
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	addAuthHeader(token, w)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in login handler")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var credentials models.Credentials

	if err = json.NewDecoder(reader).Decode(&credentials); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if err = credentials.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("invalid credentials: %s", err.Error()), http.StatusBadRequest)
		return
	}

	user, err := h.auth.Login(credentials)

	var invalidCredentialsError *services.InvalidCredentialsError
	if errors.As(err, &invalidCredentialsError) {
		http.Error(w, "incorrect login or password", http.StatusUnauthorized)
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("unexpected error in login handler")
		http.Error(w, "cannot login user", http.StatusInternalServerError)
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in login handler")
		http.Error(w, "cannot generate token", http.StatusInternalServerError)
		return
	}

	addAuthHeader(token, w)
}

func addAuthHeader(token string, w http.ResponseWriter) {
	w.Header().Set("Authorization", "Bearer "+token)
}
