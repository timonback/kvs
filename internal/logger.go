package internal

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func InitLogger(prefix string) {
	Logger = log.New(os.Stdout, prefix, log.LstdFlags)
}
