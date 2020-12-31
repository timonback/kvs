package handler

import (
	"encoding/json"
	"net/http"
)

type Leader struct {
	CurrentLeader string `json:"leader"`
	Elections     int    `json:"elections"`
}

var (
	leader = Leader{
		CurrentLeader: GetInternalId(),
		Elections:     0,
	}
)

func LeaderElection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaderMessage, _ := json.Marshal(leader)

		w.WriteHeader(http.StatusOK)
		w.Write(leaderMessage)
	})
}
