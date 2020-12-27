package store

type InmemoryService struct {
	store map[Path]Item
}

/**
Store implementation which is in-memory only
*/
func NewStoreInmemoryService() Service {
	return &InmemoryService{
		store: make(map[Path]Item),
	}
}

func (s *InmemoryService) Name() string {
	return "inmemory"
}

func (s *InmemoryService) Read(path Path) (Item, error) {
	if _, ok := s.store[path]; ok != true {
		return Item{}, NotFoundError
	}
	return s.store[path], nil
}

func (s *InmemoryService) Create(path Path, item Item) error {
	if _, ok := s.store[path]; ok == true {
		return DuplicateKeyError
	}
	s.store[path] = item
	return nil
}

func (s *InmemoryService) Update(path Path, item Item) error {
	if _, ok := s.store[path]; ok != true {
		return NotFoundError
	}
	s.store[path] = item
	return nil
}

func (s *InmemoryService) Delete(path Path) error {
	if _, ok := s.store[path]; ok != true {
		return NotFoundError
	}

	delete(s.store, path)
	return nil
}
