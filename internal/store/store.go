package store

import (
	"errors"
	"github.com/timonback/keyvaluestore/internal/store/model"
)

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
	Paths() []model.Path
	/**
	Read an item
	*/
	Read(path model.Path) (model.Item, error)
	/**
	Create a new item in the store
	Fails when an item already exists at the path
	*/
	Create(path model.Path, item model.Item) error
	/**
	Updates an existing item
	Fails when no item exists at the path
	*/
	Update(path model.Path, item model.Item) error
	/**
	Writes an item, regardless of an existing item
	*/
	Write(path model.Path, item model.Item) error
	/**
	Delete an item
	Fails when no item exists at the path
	*/
	Delete(path model.Path) error
}

var (
	NotFoundError     = errors.New("no entry found")
	DuplicateKeyError = errors.New("existing entry for key")
)
