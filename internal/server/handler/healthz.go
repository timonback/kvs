package handler

import (
	"encoding/json"
	health2 "github.com/timonback/keyvaluestore/internal/server/health"
	"net/http"
)

func Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		health := health2.GetHealthStatus()
		if health.Overall {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		writeJsonHeaders(w)
		message, _ := json.Marshal(health)
		w.Write(message)
	})
}
