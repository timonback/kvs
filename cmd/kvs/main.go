package main

import (
	"github.com/schollz/peerdiscovery"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	"github.com/timonback/keyvaluestore/internal/server"
	"time"
)

func main() {
	internal.InitLogger(false)

	arguments := arguments.ParseServerArguments()

	peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		TimeLimit: 10 * time.Second,
		Payload:   []byte(arguments.ListenPort),
		Notify: func(d peerdiscovery.Discovered) {
			arguments.Peers <- d.Address + string(d.Payload)
		},
		AllowSelf: true,
	})

	server.StartServer(&arguments)
}
