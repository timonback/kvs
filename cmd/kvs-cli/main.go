package main

import (
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	internal.InitLogger(true)

	arguments := arguments.ParseCliArguments()

	getItem(arguments)
	postItem(arguments)
	getItem(arguments)
}

func getItem(arguments arguments.Cli) {
	resp, err := http.Get(arguments.Protocol + arguments.ListenAddr + "/api/store/item")
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

func postItem(arguments arguments.Cli) {
	requestJson := `{"data": "testString"}`
	_, err := http.Post(arguments.Protocol+arguments.ListenAddr+"/api/store/item", "application/json", strings.NewReader(requestJson))
	if err != nil {
		internal.Logger.Fatal(err)
	}
}
