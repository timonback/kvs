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
	internal.InitLogger("cli")

	arguments := arguments.ParseCliArguments()

	postItem(arguments)
	getItem(arguments)

	listItem(arguments)
}

func listItem(arguments arguments.Cli) {
	resp, err := http.Get(arguments.Protocol + arguments.ListenAddr + "/api/store")
	if err != nil {
		internal.Logger.Fatal(err)
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		internal.Logger.Println("LIST " + string(body))

		resp.Body.Close()
	} else {
		internal.Logger.Println("No response")
	}
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
		internal.Logger.Println("GET " + string(body))

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
