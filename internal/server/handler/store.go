package handler

import (
	"encoding/json"
	"fmt"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"net/http"
)

type storeResponseGet struct {
	key string `json:"key"`
}

func Store() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storePath := r.URL.Path[len(context.HandlerPathStore):]

		response := storeResponseGet{
			key: storePath,
		}
		message, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(message))
	})
}
