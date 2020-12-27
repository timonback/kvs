package store

import (
	"testing"
)

var (
	key  = "key"
	path = Path(key)
	item = Item{
		Content: []byte("content"),
	}
)

func allStores() []Service {
	return []Service{NewStoreInmemoryService(), NewStoreFilesystemService("")}
}

func BenchmarkStoresGet(b *testing.B) {
	for _, store := range allStores() {
		store.Delete(path)
		store.Create(path, item)

		b.Run(store.Name(), benchmarkStoreGet(store))

		store.Delete(path)
	}
}
func benchmarkStoreGet(store Service) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.Read(path)
		}
	}
}

func BenchmarkStoresUpdate(b *testing.B) {
	for _, store := range allStores() {
		store.Delete(path)
		store.Create(path, item)

		b.Run(store.Name(), benchmarkStoreUpdate(store))

		store.Delete(path)
	}
}
func benchmarkStoreUpdate(store Service) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.Update(path, item)
		}
	}
}

func BenchmarkStoresUpdateError(b *testing.B) {
	for _, store := range allStores() {
		b.Run(store.Name(), benchmarkStoreUpdateError(store))
	}
}
func benchmarkStoreUpdateError(store Service) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.Update(path, item)
		}
	}
}

func BenchmarkStoresDeleteError(b *testing.B) {
	for _, store := range allStores() {
		b.Run(store.Name(), benchmarkStoreDeleteError(store))
	}
}
func benchmarkStoreDeleteError(store Service) func(b *testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			store.Delete(path)
		}
	}
}
