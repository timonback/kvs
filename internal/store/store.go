package store

import "errors"

type Path string

type Item struct {
	id      string
	Content string
}

/**
Interface to various store implementations
 */
type Service interface {
	/**
	Get the name of the store
	 */
	String() string
	/**
	Get all current items in the store
	Does not guarantee a refreshed look in case the store was changed externally
	 */
	Paths() []Path
	/**
	Read an item
	 */
	Read(path Path) (Item, error)
	/**
	Create a new item in the store
	Fails when an item already exists at the path
	 */
	Create(path Path, item Item) error
	/**
	Updates an existing item
	Fails when no item exists at the path
	 */
	Update(path Path, item Item) error
	/**
	Writes an item, regardless of an existing item
	 */
	Write(path Path, item Item) error
	/**
	Delete an item
	Fails when no item exists at the path
	 */
	Delete(path Path) error
}

var (
	NotFoundError     = errors.New("no entry found")
	DuplicateKeyError = errors.New("existing entry for key")
)
