package arguments

import (
	"flag"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/store"
	"os"
)

type Server struct {
	ListenAddr string
	Stop       chan os.Signal
	Store      store.Service
}

func ParseServerArguments() Server {
	arguments := Server{}
	arguments.Stop = make(chan os.Signal, 1)

	storeMode := "inmemory"
	flag.StringVar(&arguments.ListenAddr, "listen-addr", ":80", "server listen address")
	flag.StringVar(&storeMode, "store", "inmemory", "store mode. Can be inmemory, filesystem")
	flag.Parse()

	switch storeMode {
	case "inmemory":
		arguments.Store = store.NewStoreInmemoryService()
		break
	case "filesystem":
		arguments.Store = store.NewStoreFilesystemService("fs_")
		break
	default:
		internal.Logger.Println("Invalid parameter for flag store")
		os.Exit(-1)
	}

	return arguments
}
