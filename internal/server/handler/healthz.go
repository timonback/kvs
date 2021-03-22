package handler

import (
	"net/http"
	"sync/atomic"
)

const (
	HEALTHY     int32 = 0
	NON_HEALTHY int32 = 1
)

var (
	healthy int32 = NON_HEALTHY
)

func SetHealthy(healthyUpdate int32) {
	atomic.StoreInt32(&healthy, healthyUpdate)
}

func Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&healthy) == HEALTHY {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
