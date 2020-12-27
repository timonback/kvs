package arguments

import "flag"

type Server struct {
	ListenAddr string
}

func ParseServerArguments() Server {
	arguments := Server{}

	flag.StringVar(&arguments.ListenAddr, "listen-addr", ":80", "server listen address")
	flag.Parse()

	return arguments
}
