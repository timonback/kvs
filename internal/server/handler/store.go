package handler

import (
	"encoding/json"
	"errors"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	model2 "github.com/timonback/keyvaluestore/internal/store/model"
	"net/http"
	"time"
)

type storeResponse struct{}

type storeReponseList struct {
	Paths []model2.Path `json:"paths"`
}

type storeResponseGet struct {
	Key          string    `json:"key"`
	Content      string    `json:"content"`
	LastModified time.Time `json:"lastModified"`
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
			item, err := store.Read(model2.Path(storePath))
			if err != nil {
				HandleError(w, r, http.StatusNoContent, err, params)
				return
			}
			response := storeResponseGet{
				Key:          string(storePath),
				Content:      item.Content,
				LastModified: item.Time,
			}
			message, _ = json.Marshal(response)
		} else if r.Method == "POST" || r.Method == "PUT" {
			itemRequest := model.StoreRequestPost{}
			if err := MapBodyToStruct(r, &itemRequest); err != nil {
				HandleError(w, r, http.StatusBadRequest, err, params)
				return
			}
			item := model2.Item{
				Content: itemRequest.Content,
				Time:    time.Now(),
			}
			if r.URL.Query().Get(context.QueryParameterForce) != "" {
				if err := store.Write(model2.Path(storePath), item); err != nil {
					HandleError(w, r, http.StatusBadRequest, err, params)
					return
				}
			} else if r.Method == "POST" {
				if err := store.Create(model2.Path(storePath), item); err != nil {
					HandleError(w, r, http.StatusBadRequest, err, params)
					return
				}
			} else if r.Method == "PUT" {
				if err := store.Update(model2.Path(storePath), item); err != nil {
					HandleError(w, r, http.StatusBadRequest, err, params)
					return
				}
			} else {
				HandleError(w, r, http.StatusInternalServerError, errors.New("method not implemented"), nil)
			}
			message, _ = json.Marshal(storeResponse{})
		} else if r.Method == "DELETE" {
			if err := store.Delete(model2.Path(storePath)); err != nil {
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
