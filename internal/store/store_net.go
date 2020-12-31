package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/pojo"
	"github.com/timonback/keyvaluestore/internal/server/replica"
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
func NewStoreNetworkService(store Service, leader func() replica.Peer) Service {
	return &NetworkService{
		replica:   store,
		getLeader: leader,
	}
}

func (s *NetworkService) String() string {
	return "network(" + s.replica.String() + ")"
}

func (s *NetworkService) Paths() []Path {
	return s.replica.Paths()
}

func (s *NetworkService) Read(path Path) (Item, error) {
	return s.replica.Read(path)
}

func (s *NetworkService) Create(path Path, item Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := pojo.StoreRequestPost{Content: item.Content}
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
	return s.replica.Create(path, item)
}

func (s *NetworkService) Update(path Path, item Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := pojo.StoreRequestPost{Content: item.Content}
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
	return s.replica.Update(path, item)
}

func (s *NetworkService) Write(path Path, item Item) error {
	leader := replica.GetLeader()
	if leader.Id != context.GetInstanceId() {
		data := pojo.StoreRequestPost{Content: item.Content}
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
	return s.replica.Write(path, item)
}

func (s *NetworkService) Delete(path Path) error {
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
	return s.replica.Delete(path)
}
