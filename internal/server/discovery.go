package server

import (
	"github.com/schollz/peerdiscovery"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/arguments"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"io/ioutil"
	"net/http"
	"time"
)

func StartServerDiscovery(arguments *arguments.Server) {
	go peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		TimeLimit: -1,
		Delay:     3 * time.Second,
		Payload:   []byte(arguments.ListenPort),
		Notify: func(d peerdiscovery.Discovered) {
			arguments.Peers <- d.Address + ":" + string(d.Payload)
		},
		AllowSelf: true,
	})

	go func() {
		discovered := make(map[string]time.Time)

		for peer := range arguments.Peers {
			if _, ok := discovered[peer]; ok == false {
				resp, err := http.Get("http://" + peer + context.HandlerPathInternalId)
				body, _ := ioutil.ReadAll(resp.Body)
				if err == nil && len(body) == context.DiscoveryIdLength {
					internal.Logger.Printf("Discovered server at %s with id %s", peer, body)
					discovered[peer] = time.Now()
				} else {
					internal.Logger.Printf("Invalid server discovered at %s", peer)
				}
			}
		}
	}()
}