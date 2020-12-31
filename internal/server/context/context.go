package context

import "github.com/timonback/keyvaluestore/internal"

type key int

const (
	RequestIDKey    key = 0
	ApplicationJson     = "application/json"

	HandlerPathApi = "/api"

	HandlerPathStore = HandlerPathApi + "/store/"

	HandlerPathInternal                = HandlerPathApi + "/internal"
	HandlerPathInternalId              = HandlerPathInternal + "/id"
	HandlerPathInternalReplica         = HandlerPathInternal + "/replica"
	HandlerPathInternalReplicaElection = HandlerPathInternalReplica + "/leader"

	DiscoveryIdLength = 64
)

var (
	instanceId = internal.RandomString(DiscoveryIdLength)
)

func GetInstanceId() string {
	return instanceId
}
