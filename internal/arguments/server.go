package arguments

import (
	"flag"
	"os"
)

type Server struct {
	ListenAddr string
	Stop       chan os.Signal
}

func ParseServerArguments() Server {
	arguments := Server{}
	arguments.Stop = make(chan os.Signal, 1)

	flag.StringVar(&arguments.ListenAddr, "listen-addr", ":80", "server listen address")
	flag.Parse()

	return arguments
}
