package store

type FilesystemService struct{}

func NewStoreFilesystemService() Service {
	return &FilesystemService{}
}

func (s *FilesystemService) Read(path Path) (Item, error) {
	panic("implement me")
}

func (s *FilesystemService) Create(path Path, item Item) error {
	panic("implement me")
}

func (s *FilesystemService) Update(path Path, item Item) error {
	panic("implement me")
}

func (s *FilesystemService) Delete(path Path) error {
	panic("implement me")
}
