package model

import (
	"github.com/timonback/keyvaluestore/internal/store/model"
	"time"
)

type StoreRequestPost struct {
	Content string `json:"data"`
}

type StoreReponseList struct {
	Paths []model.Path `json:"paths"`
}

type StoreResponseGet struct {
	Key          string    `json:"key"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
}
