package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/belamov/ypgo-gophermart/internal/accrual/models"
	"github.com/rs/zerolog/log"
)

func (h *Handler) RegisterReward(w http.ResponseWriter, r *http.Request) {
	reader, err := getDecompressedReader(r)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register reward handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rewardCondition models.RewardCondition

	if err := json.NewDecoder(reader).Decode(&rewardCondition); err != nil {
		http.Error(w, "cannot decode json", http.StatusBadRequest)
		return
	}

	if err := rewardCondition.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := h.rewardsStorage.Exists(rewardCondition.Match)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register reward handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = h.rewardsStorage.CreateNew(rewardCondition)
	if err != nil {
		log.Error().Err(err).Msg("unexpected error in register reward handler:")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
