package handler

import (
	"encoding/json"
	"errors"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	model2 "github.com/timonback/keyvaluestore/internal/store/model"
	"github.com/timonback/keyvaluestore/internal/util"
	"net/http"
	"time"
)

type storeResponse struct{}

func Store(store store2.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storePath := r.URL.Path[len(context.HandlerPathStore):]
		params := CreateErrorParameters(ErrorParameter{
			key:   ParameterPath,
			value: storePath,
		})

		message := []byte(nil)

		if storePath == "" && r.Method == "GET" {
			paths := store.Paths()
			response := model.StoreReponseList{
				Paths: paths,
			}

			message, _ = json.Marshal(response)
		} else if r.Method == "GET" {
			item, err := store.Read(model2.Path(storePath))
			if err != nil {
				HandleError(w, r, http.StatusNoContent, err, params)
				return
			}
			response := model.StoreResponseGet{
				Key:          string(storePath),
				Content:      item.Content,
				LastModified: item.Time,
			}
			message, _ = json.Marshal(response)
		} else if r.Method == "POST" || r.Method == "PUT" {
			itemRequest := model.StoreRequestPost{}
			if err := util.MapBodyToStruct(r.Body, r.Header, &itemRequest); err != nil {
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

		writeJsonHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}
