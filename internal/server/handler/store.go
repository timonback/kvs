package handler

import (
	"encoding/json"
	"errors"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/pojo"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	"net/http"
)

type storeResponse struct{}

type storeReponseList struct {
	Paths []store2.Path `json:"paths"`
}

type storeResponseGet struct {
	Key     string `json:"key"`
	Content string `json:"content"`
}

func Store(store store2.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storePath := r.URL.Path[len(context.HandlerPathStore):]
		params := CreateErrorParameters(ErrorParameter{
			key:   PARAMETER_PATH,
			value: storePath,
		})

		message := []byte(nil)

		if storePath == "" && r.Method == "GET" {
			paths := store.Paths()
			response := storeReponseList{
				Paths: paths,
			}

			message, _ = json.Marshal(response)
		} else if r.Method == "GET" {
			item, err := store.Read(store2.Path(storePath))
			if err != nil {
				HandleError(w, r, http.StatusNoContent, err, params)
				return
			}
			response := storeResponseGet{
				Key:     string(storePath),
				Content: item.Content,
			}
			message, _ = json.Marshal(response)
		} else if r.Method == "POST" || r.Method == "PUT" {
			itemRequest := pojo.StoreRequestPost{}
			if err := MapBodyToStruct(r, &itemRequest); err != nil {
				HandleError(w, r, http.StatusBadRequest, err, params)
				return
			}
			item := store2.Item{Content: itemRequest.Content}
			if r.Method == "POST" {
				if err := store.Create(store2.Path(storePath), item); err != nil {
					HandleError(w, r, http.StatusBadRequest, err, params)
					return
				}
			} else if r.Method == "PUT" {
				if err := store.Update(store2.Path(storePath), item); err != nil {
					HandleError(w, r, http.StatusBadRequest, err, params)
					return
				}
			} else {
				HandleError(w, r, http.StatusInternalServerError, errors.New("method not implemented"), nil)
			}
			message, _ = json.Marshal(storeResponse{})
		} else if r.Method == "DELETE" {
			if err := store.Delete(store2.Path(storePath)); err != nil {
				HandleError(w, r, http.StatusBadRequest, err, params)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}
