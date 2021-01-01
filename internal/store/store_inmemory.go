package store

import "github.com/timonback/keyvaluestore/internal/store/model"

type InmemoryService struct {
	id    string
	store map[model.Path]model.Item
}

/**
Store implementation which is in-memory only
*/
func NewStoreInmemoryService(id string) Service {
	return &InmemoryService{
		id:    id,
		store: make(map[model.Path]model.Item),
	}
}

func (s *InmemoryService) String() string {
	return "inmemory" + s.id
}

func (s *InmemoryService) Paths() []model.Path {
	keys := make([]model.Path, 0, len(s.store))
	for k := range s.store {
		keys = append(keys, k)
	}
	return keys
}

func (s *InmemoryService) Read(path model.Path) (model.Item, error) {
	if _, ok := s.store[path]; ok != true {
		return model.Item{}, NotFoundError
	}
	return s.store[path], nil
}

func (s *InmemoryService) Create(path model.Path, item model.Item) error {
	if _, ok := s.store[path]; ok == true {
		return DuplicateKeyError
	}
	s.store[path] = item
	return nil
}

func (s *InmemoryService) Update(path model.Path, item model.Item) error {
	if _, ok := s.store[path]; ok != true {
		return NotFoundError
	}
	s.store[path] = item
	return nil
}

func (s *InmemoryService) Write(path model.Path, item model.Item) error {
	s.store[path] = item
	return nil
}

func (s *InmemoryService) Delete(path model.Path) error {
	if _, ok := s.store[path]; ok != true {
		return NotFoundError
	}

	delete(s.store, path)
	return nil
}
