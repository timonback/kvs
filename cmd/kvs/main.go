package main

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/cli"
	"github.com/timonback/keyvaluestore/internal/server"
)

func main() {
	internal.InitLogger(false)

	arguments := cli.ParseArguments()

	server.StartServer(&arguments)

}
