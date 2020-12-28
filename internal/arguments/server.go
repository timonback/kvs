package arguments

import (
	"flag"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/store"
	"os"
)

type Server struct {
	ListenPort string
	Stop       chan os.Signal

	Peers chan string

	Store store.Service
}

func ParseServerArguments() Server {
	arguments := Server{}

	arguments.Stop = make(chan os.Signal, 1)
	arguments.Peers = make(chan string, 1)

	storeMode := "inmemory"
	flag.StringVar(&arguments.ListenPort, "listen-port", "8080", "server listen port")
	flag.StringVar(&storeMode, "store", "filesystem", "store mode. Can be inmemory, filesystem")
	flag.Parse()

	switch storeMode {
	case "inmemory":
		arguments.Store = store.NewStoreInmemoryService("")
		break
	case "filesystem":
		pwd, _ := os.Getwd()
		folder := pwd + "/data"
		os.MkdirAll(folder, 0755)
		arguments.Store = store.NewStoreFilesystemService(folder, "")
		break
	default:
		internal.Logger.Println("Invalid parameter for flag store")
		os.Exit(-1)
	}

	return arguments
}
