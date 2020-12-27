package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timonback/keyvaluestore/internal/server/context"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	"net/http"
)

type storeRequestPost struct {
	Content interface{} `json:"data"`
}

type storeResponse struct{}

type storeResponseGet struct {
	Key     string      `json:"key"`
	Content interface{} `json:"content"`
}

func Store(store store2.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storePath := r.URL.Path[len(context.HandlerPathStore):]

		message := []byte(nil)

		if r.Method == "GET" {
			item, err := store.Read(store2.Path(storePath))
			if err != nil {
				HandleError(w, r, err, http.StatusNoContent)
				return
			}
			response := storeResponseGet{
				Key:     string(storePath),
				Content: item.Content,
			}
			message, _ = json.Marshal(response)
		} else if r.Method == "POST" || r.Method == "PUT" {
			itemRequest := storeRequestPost{}
			if err := MapBodyToStruct(r, &itemRequest); err != nil {
				HandleError(w, r, err, http.StatusBadRequest)
				return
			}
			item := store2.Item{Content: itemRequest.Content}
			if r.Method == "POST" {
				if err := store.Create(store2.Path(storePath), item); err != nil {
					HandleError(w, r, err, http.StatusBadRequest)
					return
				}
			} else if r.Method == "PUT" {
				if err := store.Update(store2.Path(storePath), item); err != nil {
					HandleError(w, r, err, http.StatusBadRequest)
					return
				}
			} else {
				HandleError(w, r, errors.New("method not implemented"), http.StatusInternalServerError)
			}
			message, _ = json.Marshal(storeResponse{})
		} else if r.Method == "DELETE" {
			if err := store.Delete(store2.Path(storePath)); err != nil {
				HandleError(w, r, err, http.StatusBadRequest)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, string(message))
	})
}
