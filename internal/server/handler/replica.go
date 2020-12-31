package handler

import (
	"encoding/json"
	"github.com/timonback/keyvaluestore/internal/server/replica"
	"net/http"
)

func LeaderElection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetLeader())

		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}
func Peers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetOnlinePeers())

		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}
