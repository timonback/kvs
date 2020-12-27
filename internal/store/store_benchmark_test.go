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
	inmemory1 := NewStoreInmemoryService("1")
	inmemory2 := NewStoreInmemoryService("2")
	inmemory3 := NewStoreInmemoryService("3")
	inmemory4 := NewStoreInmemoryService("4")
	inmemory5 := NewStoreInmemoryService("5")
	filesystem1 := NewStoreFilesystemService("")
	filesystem2 := NewStoreFilesystemService("1_")
	replica := NewStoreReplicaService(inmemory1, inmemory2, inmemory3, inmemory4, inmemory5)
	replica2 := NewStoreReplicaService(filesystem1, filesystem2)
	return []Service{inmemory1, filesystem1, replica, replica2}
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
