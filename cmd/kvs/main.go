package main

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	"github.com/timonback/keyvaluestore/internal/server"
	store2 "github.com/timonback/keyvaluestore/internal/store"
)

func main() {
	internal.InitLogger(false)

	arguments := arguments.ParseServerArguments()
	store := store2.NewStoreInmemoryService()

	server.StartServer(&arguments, store)

}
