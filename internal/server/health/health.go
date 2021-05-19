package health

import (
	"sync"
	"sync/atomic"
)

const (
	HEALTHY        int32 = 0b00000000
	SERVER_STATUS  int32 = 0b00000001
	REPLICA_STATUS int32 = 0b00000010
)

var (
	healthy = HEALTHY
	lock    sync.Mutex
)

type Status struct {
	Overall bool `json:"global"`
	Server  bool `json:"server"`
	Replica bool `json:"replica"`
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
		Server:  health&SERVER_STATUS == 0,
		Replica: health&REPLICA_STATUS == 0,
	}
}
