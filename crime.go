package main

import "net/http"

func (cfg *apiConfig) handlerCrimeStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	stats, err := cfg.db.GetAllCrimeStats(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch crime stats", err)
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}
