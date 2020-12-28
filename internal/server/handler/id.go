package handler

import (
	"fmt"
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
		fmt.Fprint(w, id)
	})
}
