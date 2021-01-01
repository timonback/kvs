package store

import (
	"bufio"
	"encoding/base64"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/store/model"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	timeFormat = time.RFC3339
)

type FilesystemService struct {
	folder     string
	filePrefix string
}

/**
Store implementation which uses the filesystem for persistent storage

File-format: Header\n\nContent
Header can be n lines, terminated by \n. A header may never be an empty line
Content can be any bytes till the end of the file
*/
func NewStoreFilesystemService(folder string, filePrefix string) Service {
	return &FilesystemService{
		folder:     folder,
		filePrefix: filePrefix,
	}
}

func (s *FilesystemService) String() string {
	return "filesystem"
}

func (s *FilesystemService) Paths() []model.Path {
	files, _ := ioutil.ReadDir(s.folder)

	keys := make([]model.Path, 0, len(files))
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

func (s *FilesystemService) filenameToPath(file os.FileInfo) model.Path {
	prefix := s.filePrefix + "_fs_"
	if strings.HasPrefix(file.Name(), prefix) {
		base64Name := file.Name()[len(prefix):]
		key, err := base64.StdEncoding.DecodeString(base64Name)
		if err == nil {
			return model.Path(key)
		}
	}
	return ""
}
func (s *FilesystemService) pathToFilename(path model.Path) string {
	return s.folder + "/" + s.filePrefix + "_fs_" + base64.StdEncoding.EncodeToString([]byte(path))
}

func (s *FilesystemService) fileExists(path model.Path) bool {
	_, err := os.Stat(s.pathToFilename(path))
	return !os.IsNotExist(err)
}

func (s *FilesystemService) Read(path model.Path) (model.Item, error) {
	filename := s.pathToFilename(path)
	file, err := os.Open(filename)
	if err != nil {
		return model.Item{}, NotFoundError
	}
	defer file.Close()

	item := model.Item{}

	reader := bufio.NewReader(file)
	for {
		headerLine, err := reader.ReadString('\n')
		if err != nil {
			panic("Read invalid filesystem file for path " + string(path))
		}

		if headerLine == "\n" {
			// Empty header line. Content is coming next
			break
		}

		header := strings.Split(headerLine, ":")
		headerKey := header[0]
		headerValue := headerLine[len(header[0])+1 : len(headerLine)-1]
		switch headerKey {
		case "time":
			timeVal, err := time.Parse(timeFormat, headerValue)
			if err == nil {
				item.Time = timeVal
			} else {
				internal.Logger.Println("Filesystem store file has invalid time header " + string(path))
			}
			break
		default:
			internal.Logger.Printf("Found invalid header line in filesystem: " + headerLine)
		}

	}

	content, _ := ioutil.ReadAll(reader)
	item.Content = string(content)

	return item, nil
}

func (s *FilesystemService) Create(path model.Path, item model.Item) error {
	if s.fileExists(path) {
		return DuplicateKeyError
	}
	return s.writeContentToDisk(path, item)
}

func (s *FilesystemService) Update(path model.Path, item model.Item) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return s.writeContentToDisk(path, item)
}

func (s *FilesystemService) Write(path model.Path, item model.Item) error {
	return s.writeContentToDisk(path, item)
}

func (s *FilesystemService) Delete(path model.Path) error {
	if !s.fileExists(path) {
		return NotFoundError
	}
	return os.Remove(s.pathToFilename(path))
}

func (s *FilesystemService) writeContentToDisk(path model.Path, item model.Item) error {
	return ioutil.WriteFile(s.pathToFilename(path), []byte("time:"+time.Now().Format(timeFormat)+"\n\n"+item.Content), 0744)
}
