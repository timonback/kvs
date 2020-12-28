package arguments

import "flag"

type Cli struct {
	ListenAddr string
	Protocol   string
}

func ParseCliArguments() Cli {
	arguments := Cli{}

	flag.StringVar(&arguments.ListenAddr, "listen-addr", "localhost:8080", "server listen address")
	flag.StringVar(&arguments.Protocol, "protocol", "http://", "server listen protocol")
	flag.Parse()

	return arguments
}
