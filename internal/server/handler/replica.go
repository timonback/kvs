package handler

import (
	"encoding/json"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	"github.com/timonback/keyvaluestore/internal/server/replica"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	model2 "github.com/timonback/keyvaluestore/internal/store/model"
	"github.com/timonback/keyvaluestore/internal/util"
	"net/http"
	"strconv"
)

func PeerStatus(store store2.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := model.StoreResponseReplicaStatus{
			Id:             context.GetInstanceId(),
			Uptime:         context.GetUpTime(),
			LogBookEntries: 0,
			IsLeader:       context.GetInstanceId() == replica.GetLeader().Id,
		}
		if logbook, err := store.Read(context.LogBookEntryStorePath); err == nil {
			logBookEntries, _ := strconv.Atoi(logbook.Content)
			status.LogBookEntries = logBookEntries
		}
		message, _ := json.Marshal(status)

		writeJsonHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}

func LeaderElection() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetLeader())

		writeJsonHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}

func Peers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		message, _ := json.Marshal(replica.GetOnlinePeers(true))

		writeJsonHeaders(w)
		w.WriteHeader(http.StatusOK)
		w.Write(message)
	})
}

func StoreSync(store *store2.NetworkService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		itemRequest := model.StoreRequestSync{}
		if err := util.MapBodyToStruct(r.Body, r.Header, &itemRequest); err != nil {
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
			store.UpdateLogbookEntry()
		}

		writeJsonHeaders(w)
		w.WriteHeader(http.StatusNoContent)
	})
}
