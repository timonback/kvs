package replica

import (
	"fmt"
	"github.com/schollz/peerdiscovery"
	"github.com/timonback/keyvaluestore/internal"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/model"
	"github.com/timonback/keyvaluestore/internal/server/health"
	"github.com/timonback/keyvaluestore/internal/util"
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

var (
	state    = model.Election
	peers    = make(map[string]Peer)
	leaderId = ""
)

func StartServerDiscovery(listenPort int) chan model.StoreResponseReplicaStatus {
	internal.Logger.Println("Discovery is starting...")
	health.SetUnhealthy(health.REPLICA_STATUS)

	discoveredPeers := make(chan string, 10)
	newLeader := make(chan model.StoreResponseReplicaStatus, 1)

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
	goCheckPeersHealth(newLeader)

	return newLeader
}

func goVerifyNewPeers(discoveredPeers chan string) {
	go func() {
		for peerAddress := range discoveredPeers {
			if _, ok := peers[peerAddress]; ok == false {
				if available, status := IsPeerAvailable(peerAddress); available {
					internal.Logger.Printf("Discovered server at %s with id %s", peerAddress, status.Id)
					now := time.Now()
					peers[peerAddress] = Peer{
						Id:           status.Id,
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

func goCheckPeersHealth(newLeader chan model.StoreResponseReplicaStatus) {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				leadershipCheck(newLeader)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func leadershipCheck(newLeader chan model.StoreResponseReplicaStatus) {
	// TODO: Improve leadership check by syncing between all instances (agree on leadership at the same time)
	leader := model.StoreResponseReplicaStatus{}
	for _, peer := range peers {
		if available, peerStatus := IsPeerAvailable(peer.Address); !available || peerStatus.Id != peer.Id {
			internal.Logger.Printf("Removing unavailable/restarted peer %v", peer)
			delete(peers, peer.Address)
			if peer.Id == leaderId {
				internal.Logger.Printf("Removed peer was the leader. Restarting election")
				leaderId = ""
				state = model.Election
				health.SetUnhealthy(health.REPLICA_STATUS)
			}
		} else if leader.Id == "" || leader.LogBookEntries < peerStatus.LogBookEntries || (leader.LogBookEntries == peerStatus.LogBookEntries && leader.Uptime.After(peerStatus.Uptime)) {
			if peerStatus.State == model.Primary || peerStatus.State == model.Election {
				leader = peerStatus
			}
		}
	}
	if leader.Id != "" && leaderId != leader.Id {
		leaderId = leader.Id
		if leader.Id != context.GetInstanceId() {
			internal.Logger.Printf(fmt.Sprintf("Leadership lost. New leader is at %s", GetLeader().Address))
			state = model.Secondary
		} else {
			internal.Logger.Println("Leadership gained")
			state = model.Primary
		}

		health.SetHealthy(health.REPLICA_STATUS)
		newLeader <- leader
	}
}

func IsPeerAvailable(peerAddress string) (bool, model.StoreResponseReplicaStatus) {
	response := model.StoreResponseReplicaStatus{}

	r, err := http.Get("http://" + peerAddress + context.HandlerPathInternalReplicaStatus)
	if err == nil {
		err := util.MapBodyToStruct(r.Body, r.Header, &response)
		if err == nil && len(response.Id) == context.DiscoveryIdLength {
			return true, response
		}
	}
	return false, response
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

func GetReplicaState() model.StoreReplicaState {
	return state
}
