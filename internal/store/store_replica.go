package store

import (
	"github.com/timonback/keyvaluestore/internal/store/model"
	"sync"
)

type ReplicaService struct {
	replicas []Service
}

/**
Store implementation which coordinates multiple local store implementations
Commands are simply forwarded to the stores.
The first store functions as the master. Create/Update/Delete commands are only forwarded to the other stores without
further error checks
*/
func NewStoreReplicaService(stores ...Service) Service {
	if len(stores) == 0 {
		panic("the replica store requires at least one store")
	}
	return &ReplicaService{
		replicas: stores,
	}
}

func (s *ReplicaService) String() string {
	name := "replicas(" + s.replicas[0].String()
	for _, replica := range s.replicas[1:] {
		name += "," + replica.String()
	}
	name += ")"

	return name
}

func (s *ReplicaService) Paths() []model.Path {
	return s.replicas[0].Paths()
}

func (s *ReplicaService) Read(path model.Path) (model.Item, error) {
	return s.replicas[0].Read(path)
}

func (s *ReplicaService) Create(path model.Path, item model.Item) error {
	err := s.replicas[0].Create(path, item)
	if err == nil {
		wg := sync.WaitGroup{}
		for _, replica := range s.replicas[1:] {
			wg.Add(1)
			go func(rep Service) {
				rep.Create(path, item)
				wg.Done()
			}(replica)
		}
		wg.Wait()
	}
	return err
}

func (s *ReplicaService) Update(path model.Path, item model.Item) error {
	err := s.replicas[0].Update(path, item)
	if err == nil {
		wg := sync.WaitGroup{}
		for _, replica := range s.replicas[1:] {
			wg.Add(1)
			go func(rep Service) {
				rep.Update(path, item)
				wg.Done()
			}(replica)
		}
		wg.Wait()
	}
	return err
}

func (s *ReplicaService) Write(path model.Path, item model.Item) error {
	err := s.replicas[0].Write(path, item)
	if err == nil {
		wg := sync.WaitGroup{}
		for _, replica := range s.replicas[1:] {
			wg.Add(1)
			go func(rep Service) {
				rep.Write(path, item)
				wg.Done()
			}(replica)
		}
		wg.Wait()
	}
	return err
}

func (s *ReplicaService) Delete(path model.Path) error {
	err := s.replicas[0].Delete(path)
	if err == nil {
		wg := sync.WaitGroup{}
		for _, replica := range s.replicas[1:] {
			wg.Add(1)
			go func(rep Service) {
				rep.Delete(path)
				wg.Done()
			}(replica)
		}
		wg.Wait()
	}
	return err
}
