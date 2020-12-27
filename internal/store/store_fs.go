package store

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"
)

type FilesystemService struct {
	folder     string
	filePrefix string
}

/**
Store implementation which uses the filesystem for persistent storage
*/
func NewStoreFilesystemService(folder string, filePrefix string) Service {
	return &FilesystemService{
		folder:     folder,
		filePrefix: filePrefix,
	}
}

func (s *FilesystemService) Name() string {
	return "filesystem"
}

func (s *FilesystemService) Paths() []Path {
	files, _ := ioutil.ReadDir(s.folder)

	keys := make([]Path, 0, len(files))
	for _, f := range files {
		if !f.IsDir() {
			key := s.filenameToPath(f)
			if key != "" {
				keys = append(keys, key)
			}
		}
	}

	return keys
}

func (s *FilesystemService) filenameToPath(file os.FileInfo) Path {
	prefix := s.filePrefix + "_fs_"
	if strings.HasPrefix(file.Name(), prefix) {
		base64Name := file.Name()[len(prefix):]
		key, err := base64.StdEncoding.DecodeString(base64Name)
		if err == nil {
			return Path(key)
		}
	}
	return ""
}
func (s *FilesystemService) pathToFilename(path Path) string {
	return s.folder + "/" + s.filePrefix + "_fs_" + base64.StdEncoding.EncodeToString([]byte(path))
}

func (s *FilesystemService) fileExists(path Path) bool {
	_, err := os.Stat(s.pathToFilename(path))
	return !os.IsNotExist(err)
}

func (s *FilesystemService) Read(path Path) (Item, error) {
	content, err := ioutil.ReadFile(s.pathToFilename(path))
	if err != nil {
		return Item{}, NotFoundError
	}
	return Item{
		Content: string(content),
	}, nil
}

func (s *FilesystemService) Create(path Path, item Item) error {
	if s.fileExists(path) {
		return DuplicateKeyError
	}
	return ioutil.WriteFile(s.pathToFilename(path), []byte(item.Content), 0744)
}

func (s *FilesystemService) Update(path Path, item Item) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return ioutil.WriteFile(s.pathToFilename(path), []byte(item.Content), 0744)
}

func (s *FilesystemService) Delete(path Path) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return os.Remove(s.pathToFilename(path))
}
