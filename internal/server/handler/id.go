package handler

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"net/http"
)

var (
	id = internal.RandomString(context.DiscoveryIdLength)
)

func InternalId() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(id))
	})
}

func GetInternalId() string {
	return id
}
