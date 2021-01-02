package arguments

import (
	"flag"
	"github.com/timonback/keyvaluestore/internal/store"
	"os"
)

type Server struct {
	ListenPort int
	Stop       chan os.Signal

	Store        store.Service
	NetworkStore *store.NetworkService // possibly duplicated. Is type casting Store -> NetworkStore possible (currently in server.go)?
}

func ParseServerArguments() Server {
	arguments := Server{}

	arguments.Stop = make(chan os.Signal, 1)

	flag.IntVar(&arguments.ListenPort, "listen-port", 8080, "server listen port")
	withFilesystemStore := flag.Bool("filesystem", true, "use the filesystem store")
	withNetworkStore := flag.Bool("network", true, "sync kvs with multiple instances in the same network")
	flag.Parse()

	arguments.Store = store.NewStoreInmemoryService("")
	if *withFilesystemStore {
		pwd, _ := os.Getwd()
		folder := pwd + "/data"
		_ = os.MkdirAll(folder, 0755)
		filesystemStore := store.NewStoreFilesystemService(folder, "")

		// init in memory store with current content from filesystem
		for _, path := range filesystemStore.Paths() {
			item, err := filesystemStore.Read(path)
			if err == nil {
				arguments.Store.Write(path, item)
			}
		}

		arguments.Store = store.NewStoreReplicaService(arguments.Store, filesystemStore)
	}
	if *withNetworkStore {
		arguments.NetworkStore = store.NewStoreNetworkService(arguments.Store, arguments.ListenPort)
		arguments.Store = arguments.NetworkStore
	}

	return arguments
}
