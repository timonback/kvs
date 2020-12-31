package main

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	"github.com/timonback/keyvaluestore/internal/server"
)

func main() {
	internal.InitLogger("server")

	arguments := arguments.ParseServerArguments()

	server.StartServer(&arguments)
}
