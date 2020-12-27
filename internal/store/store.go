package store

import "errors"

type Path string

type Item struct {
	id      string
	Content string
}

type Service interface {
	Name() string
	Paths() []Path
	Read(path Path) (Item, error)
	Create(path Path, item Item) error
	Update(path Path, item Item) error
	Delete(path Path) error
}

var (
	NotFoundError     = errors.New("no entry found")
	DuplicateKeyError = errors.New("existing entry for key")
)
