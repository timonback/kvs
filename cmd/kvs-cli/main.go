package main

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/cli"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	internal.InitLogger(true)

	arguments := cli.ParseArguments()

	resp, err := http.Get(arguments.Protocol + arguments.ListenAddr)
	if err != nil {
		internal.Logger.Fatal(err)
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		internal.Logger.Println(string(body))

		resp.Body.Close()
	} else {
		internal.Logger.Println("No response")
	}
}
