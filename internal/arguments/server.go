package arguments

import (
	"flag"
	"github.com/timonback/keyvaluestore/internal/server/replica"
	"github.com/timonback/keyvaluestore/internal/store"
	"os"
)

type Server struct {
	ListenPort int
	Stop       chan os.Signal

	DiscoveredPeers chan string

	Store store.Service
}

func ParseServerArguments() Server {
	arguments := Server{}

	arguments.Stop = make(chan os.Signal, 1)
	arguments.DiscoveredPeers = make(chan string, 10)

	flag.IntVar(&arguments.ListenPort, "listen-port", 8080, "server listen port")
	withFilesystemStore := flag.Bool("filesystem", true, "use the filesystem store")
	withNetworkStore := flag.Bool("network", true, "sync kvs with multiple instances in the same network")
	flag.Parse()

	arguments.Store = store.NewStoreInmemoryService("")
	if *withFilesystemStore {
		pwd, _ := os.Getwd()
		folder := pwd + "/data"
		os.MkdirAll(folder, 0755)
		arguments.Store = store.NewStoreReplicaService(arguments.Store, store.NewStoreFilesystemService(folder, ""))
	}
	if *withNetworkStore {
		arguments.Store = store.NewStoreNetworkService(arguments.Store, replica.GetLeader)
	}

	return arguments
}
