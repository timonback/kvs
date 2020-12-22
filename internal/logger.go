package internal

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger(isCli bool) {
	if isCli {
		Logger = log.New(os.Stdout, "cli: ", log.LstdFlags)
	} else {
		Logger = log.New(os.Stdout, "server: ", log.LstdFlags)
	}
}
