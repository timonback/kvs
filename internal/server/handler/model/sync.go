package model

import (
	"github.com/timonback/keyvaluestore/internal/store/model"
	"time"
)

type StoreSyncMode string

const (
	StoreSyncModeDelete StoreSyncMode = "DELETE"
	StoreSyncModeWrite  StoreSyncMode = "WRITE"
)

type StoreSync struct {
	Content      string        `json:"data"`
	Path         model.Path    `json:"path"`
	LastModified time.Time     `json:"lastModified"`
	Mode         StoreSyncMode `json:"mode"`
}

type StoreRequestSync struct {
	Commands []StoreSync `json:"commands"`
}

type StoreReplicaState int

const (
	Primary StoreReplicaState = iota
	Secondary
	Election
)

type StoreResponseReplicaStatus struct {
	Id             string            `json:"id"`
	Uptime         time.Time         `json:"uptime"`
	LogBookEntries int               `json:"logbookEntries"`
	State          StoreReplicaState `json:"state"`
}
