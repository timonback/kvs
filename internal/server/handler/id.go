package handler

import (
	"github.com/timonback/keyvaluestore/internal/server/context"
	"net/http"
)

func InternalId() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(context.GetInstanceId()))
	})
}
