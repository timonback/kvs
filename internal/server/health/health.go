package health

import (
	"sync"
	"sync/atomic"
)

const (
	HEALTHY            int32 = 0b00000000
	NON_HEALTHY_SERVER int32 = 0b00000001
	NON_HEALTHY_LEADER int32 = 0b00000010
)

var (
	healthy = HEALTHY
	lock    sync.Mutex
)

type Status struct {
	Overall bool `json:"global"`
	Server  bool `json:"server"`
	Leader  bool `json:"leader"`
}

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

func GetHealthStatus() Status {
	health := atomic.LoadInt32(&healthy)
	return Status{
		Overall: health == HEALTHY,
		Server:  health&NON_HEALTHY_SERVER == 0,
		Leader:  health&NON_HEALTHY_LEADER == 0,
	}
}
