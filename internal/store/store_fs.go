package store

import (
	"io/ioutil"
	"os"
	"strings"
)

type FilesystemService struct {
	pathPrefix string
}

/**
Store implementation which uses the filesystem for persistent storage
*/
func NewStoreFilesystemService(pathPrefix string) Service {
	return &FilesystemService{
		pathPrefix: pathPrefix,
	}
}

func (s *FilesystemService) Name() string {
	return "filesystem"
}

func (s *FilesystemService) pathToFilename(path Path) string {
	wd, _ := os.Getwd()

	return wd + "/" + s.pathPrefix + "_fs_" + strings.ReplaceAll(string(path), "/", "_")
}

func (s *FilesystemService) fileExists(path Path) bool {
	_, err := os.Stat(s.pathToFilename(path))
	return os.IsNotExist(err)
}

func (s *FilesystemService) Read(path Path) (Item, error) {
	content, err := ioutil.ReadFile(s.pathToFilename(path))
	if err != nil {
		return Item{}, NotFoundError
	}
	return Item{
		Content: content,
	}, nil
}

func (s *FilesystemService) Create(path Path, item Item) error {
	if s.fileExists(path) {
		return DuplicateKeyError
	}
	return ioutil.WriteFile(s.pathToFilename(path), item.Content, 0)
}

func (s *FilesystemService) Update(path Path, item Item) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return ioutil.WriteFile(s.pathToFilename(path), item.Content, 0)
}

func (s *FilesystemService) Delete(path Path) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return os.Remove(s.pathToFilename(path))
}
