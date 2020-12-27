package store

import (
	"io/ioutil"
	"os"
	"strings"
)

type FilesystemService struct{}

func NewStoreFilesystemService() Service {
	return &FilesystemService{}
}

func pathToFilename(path Path) string {
	wd, _ := os.Getwd()

	return wd + "/store_fs_" + strings.ReplaceAll(string(path), "/", "_")
}

func fileExists(path Path) bool {
	_, err := os.Stat(pathToFilename(path))
	return os.IsNotExist(err)
}

func (s *FilesystemService) Name() string {
	return "filesystem"
}

func (s *FilesystemService) Read(path Path) (Item, error) {
	content, err := ioutil.ReadFile(pathToFilename(path))
	if err != nil {
		return Item{}, NotFoundError
	}
	return Item{
		Content: content,
	}, nil
}

func (s *FilesystemService) Create(path Path, item Item) error {
	if fileExists(path) {
		return DuplicateKeyError
	}
	return ioutil.WriteFile(pathToFilename(path), item.Content, 0)
}

func (s *FilesystemService) Update(path Path, item Item) error {
	if !fileExists(path) {
		return NotFoundError
	}
	return ioutil.WriteFile(pathToFilename(path), item.Content, 0)
}

func (s *FilesystemService) Delete(path Path) error {
	if !fileExists(path) {
		return NotFoundError
	}
	return os.Remove(pathToFilename(path))
}
