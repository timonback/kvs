package cli

import "flag"

type Arguments struct {
	ListenAddr string
	Protocol   string
}

func ParseArguments() Arguments {
	arguments := Arguments{}

	flag.StringVar(&arguments.ListenAddr, "listen-addr", ":80", "server listen address")
	flag.StringVar(&arguments.Protocol, "protocol", "http://", "server listen protocol")
	flag.Parse()

	return arguments
}
