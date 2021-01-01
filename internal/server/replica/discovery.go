package replica

import (
	"github.com/schollz/peerdiscovery"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Peer struct {
	Id           string    `json:"id"`
	Address      string    `json:"address"`
	DiscoveredAt time.Time `json:"discoveredAt"`
	LastOnline   time.Time `json:"lastOnline"`
	IsOnline     bool
}

type Leader struct {
	CurrentLeader string `json:"leader"`
	NumElections  int    `json:"elections"`
}

var (
	peers    = make(map[string]Peer)
	leaderId = ""
)

func StartServerDiscovery(listenPort int, discoveredPeers chan string) {
	internal.Logger.Println("Discovery is starting...")

	go peerdiscovery.Discover(peerdiscovery.Settings{
		Limit:     -1,
		TimeLimit: -1,
		Delay:     3 * time.Second,
		Payload:   []byte(strconv.Itoa(listenPort)),
		Notify: func(d peerdiscovery.Discovered) {
			discoveredPeers <- d.Address + ":" + string(d.Payload)
		},
		AllowSelf: true,
	})

	goVerifyNewPeers(discoveredPeers)
	goCheckPeersHealth()
}

func goVerifyNewPeers(discoveredPeers chan string) {
	go func() {
		for peerAddress := range discoveredPeers {
			if _, ok := peers[peerAddress]; ok == false {
				if available, id := IsPeerAvailable(peerAddress); available {
					internal.Logger.Printf("Discovered server at %s with id %s", peerAddress, id)
					now := time.Now()
					peers[peerAddress] = Peer{
						Id:           string(id),
						Address:      peerAddress,
						DiscoveredAt: now,
						LastOnline:   now,
						IsOnline:     true,
					}
				} else {
					internal.Logger.Printf("Invalid server discovered at %s", peerAddress)
				}
			}
		}
	}()
}

func goCheckPeersHealth() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				leadershipCheck()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func leadershipCheck() {
	newLeaderId := "zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz"
	for _, peer := range peers {
		if available, id := IsPeerAvailable(peer.Address); !available || string(id) != peer.Id {
			internal.Logger.Printf("Removing unavailable/restarted peer %v", peer)
			delete(peers, peer.Address)
		} else if peer.Id < newLeaderId {
			// TODO: Improve leader election
			newLeaderId = peer.Id
		}
	}
	if leaderId != newLeaderId {
		if newLeaderId != context.GetInstanceId() {
			internal.Logger.Printf("Leadership lost")
		} else {
			internal.Logger.Println("Leadership gained")
		}
	}
	leaderId = newLeaderId
}

func IsPeerAvailable(peerAddress string) (bool, []byte) {
	resp, err := http.Get("http://" + peerAddress + context.HandlerPathInternalId)
	if err == nil {
		body, _ := ioutil.ReadAll(resp.Body)
		if len(body) == context.DiscoveryIdLength {
			return true, body
		}
	}
	return false, nil
}

func GetOnlinePeers(includeMyself bool) []Peer {
	p := make([]Peer, 0, len(peers))

	for _, peer := range peers {
		if peer.IsOnline {
			if includeMyself || peer.Id != context.GetInstanceId() {
				p = append(p, peer)
			}
		}
	}

	return p
}

func GetLeader() Peer {
	for _, peer := range peers {
		if leaderId == peer.Id {
			return peer
		}
	}
	return Peer{}
}
