package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/belamov/ypgo-gophermart/internal/accrual/storage"
	"github.com/rs/zerolog/log"
)

func (h *Handler) RegisterReward(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register reward handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rewardCondition models.Reward

	if err := json.NewDecoder(reader).Decode(&rewardCondition); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if err := rewardCondition.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.rewardsStorage.Save(rewardCondition)
	var notUniqueError *storage.NotUniqueError
	if errors.As(err, &notUniqueError) {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register reward handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
