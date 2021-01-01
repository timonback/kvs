package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	"github.com/timonback/keyvaluestore/internal/server/replica"
	model2 "github.com/timonback/keyvaluestore/internal/store/model"
	"io/ioutil"
	"net/http"
	"strings"
)

type NetworkService struct {
	replica   Service
	getLeader func() replica.Peer
}

/**
Store implementation which internally uses the internal store for read and forwards manipulation actions to the getLeader
Local store is only eventually consistent. Write/Delete actions are not waited for till completion
*/
func NewStoreNetworkService(store Service, leader func() replica.Peer) *NetworkService {
	return &NetworkService{
		replica:   store,
		getLeader: leader,
	}
}

func (s *NetworkService) String() string {
	return "network(" + s.replica.String() + ")"
}

func (s *NetworkService) Paths() []model2.Path {
	return s.replica.Paths()
}

func (s *NetworkService) Read(path model2.Path) (model2.Item, error) {
	return s.replica.Read(path)
}

func (s *NetworkService) Create(path model2.Path, item model2.Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := model.StoreRequestPost{Content: item.Content}
		body, _ := json.Marshal(data)
		resp, err := http.Post("http://"+leader.Address+context.HandlerPathStore+string(path), context.ApplicationJson, strings.NewReader(string(body)))
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				response, _ := ioutil.ReadAll(resp.Body)
				return errors.New(fmt.Sprintf("Unexpected response during creation %s (%d)", response, resp.StatusCode))
			}
		}
		return err
	}

	err := s.replica.Create(path, item)

	if err == nil {
		s.syncStoreChanges(path, model.StoreSyncModeWrite, item)
	}
	return err
}

func (s *NetworkService) syncStoreChanges(path model2.Path, mode model.StoreSyncMode, item model2.Item) {
	data := model.StoreRequestSync{
		Commands: []model.StoreSync{},
	}
	data.Commands = append(data.Commands, model.StoreSync{
		Content:      item.Content,
		Path:         path,
		LastModified: item.Time,
		Mode:         mode,
	})
	body, _ := json.Marshal(data)
	for _, peer := range replica.GetOnlinePeers(false) {
		http.Post("http://"+peer.Address+context.HandlerPathInternalReplicaSync, context.ApplicationJson, strings.NewReader(string(body)))
	}
}

func (s *NetworkService) Update(path model2.Path, item model2.Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := model.StoreRequestPost{Content: item.Content}
		body, _ := json.Marshal(data)
		resp, err := http.Post("http://"+leader.Address+context.HandlerPathStore+string(path), context.ApplicationJson, strings.NewReader(string(body)))
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				response, _ := ioutil.ReadAll(resp.Body)
				return errors.New(fmt.Sprintf("Unexpected response during creation %s (%d)", response, resp.StatusCode))
			}
		}
		return err
	}

	err := s.replica.Update(path, item)

	if err == nil {
		s.syncStoreChanges(path, model.StoreSyncModeWrite, item)
	}
	return err
}

func (s *NetworkService) Write(path model2.Path, item model2.Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := model.StoreRequestPost{Content: item.Content}
		body, _ := json.Marshal(data)
		resp, err := http.Post("http://"+leader.Address+context.HandlerPathStore+string(path), context.ApplicationJson, strings.NewReader(string(body)))
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				response, _ := ioutil.ReadAll(resp.Body)
				return errors.New(fmt.Sprintf("Unexpected response during creation %s (%d)", response, resp.StatusCode))
			}
		}
		return err
	}
	err := s.replica.Write(path, item)

	if err == nil {
		s.syncStoreChanges(path, model.StoreSyncModeWrite, item)
	}
	return err
}

func (s *NetworkService) Delete(path model2.Path) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		req, _ := http.NewRequest("DELETE", "http://"+leader.Address+context.HandlerPathStore+string(path), nil)
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			if resp.StatusCode != http.StatusOK {
				response, _ := ioutil.ReadAll(resp.Body)
				return errors.New(fmt.Sprintf("Unexpected response during creation %s (%d)", response, resp.StatusCode))
			}
		}
		return err
	}
	err := s.replica.Delete(path)

	if err == nil {
		s.syncStoreChanges(path, model.StoreSyncModeDelete, model2.Item{})
	}
	return err
}

func (s *NetworkService) GetUnderlyingService() Service {
	return s.replica
}
