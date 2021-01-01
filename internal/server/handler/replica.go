package handler

import (
	"encoding/json"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	"github.com/timonback/keyvaluestore/internal/server/replica"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	model2 "github.com/timonback/keyvaluestore/internal/store/model"
	"net/http"
)

func LeaderElection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetLeader())

		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}

func Peers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetOnlinePeers(true))

		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}

func StoreSync(store *store2.NetworkService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemRequest := model.StoreRequestSync{}
		if err := MapBodyToStruct(r, &itemRequest); err != nil {
			HandleError(w, r, http.StatusBadRequest, err, nil)
			return
		}

		for _, command := range itemRequest.Commands {
			if command.Mode == model.StoreSyncModeDelete {
				store.GetUnderlyingService().Delete(command.Path)
			} else if command.Mode == model.StoreSyncModeWrite {
				item := model2.Item{Content: command.Content, Time: command.LastModified}
				store.GetUnderlyingService().Write(command.Path, item)
			} else {
				panic("Invalid store sync mode " + command.Mode)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
