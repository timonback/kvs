package handler

import (
	"net/http"
	"sync"
	"sync/atomic"
)

const (
	HEALTHY              int32 = 0b00000000
	NON_HEALTHY_SERVER   int32 = 0b00000001
	NON_HEALTHY_ELECTION int32 = 0b00000010
)

var (
	healthy = HEALTHY
	lock    sync.Mutex
)

func SetHealthy(healthyUpdate int32) {
	lock.Lock()
	defer lock.Unlock()

	health := atomic.LoadInt32(&healthy)
	health = health & ^healthyUpdate
	atomic.StoreInt32(&healthy, health)
}

func SetUnhealthy(healthyUpdate int32) {
	lock.Lock()
	defer lock.Unlock()

	health := atomic.LoadInt32(&healthy)
	health = health | healthyUpdate
	atomic.StoreInt32(&healthy, health)
}

func IsHealth() bool {
	return atomic.LoadInt32(&healthy) == HEALTHY
}

func Healthz() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if IsHealth() {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
