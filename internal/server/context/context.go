package context

import (
	"github.com/timonback/keyvaluestore/internal/store/model"
	"github.com/timonback/keyvaluestore/internal/util"
)

type key int

const (
	RequestIDKey    key = 0
	ApplicationJson     = "application/json"

	HandlerPathApi = "/api"

	HandlerPathStore = HandlerPathApi + "/store/"

	HandlerPathInternal                = HandlerPathApi + "/internal"
	HandlerPathInternalId              = HandlerPathInternal + "/id"
	HandlerPathInternalReplica         = HandlerPathInternal + "/replica"
	HandlerPathInternalReplicaStatus   = HandlerPathInternalReplica + "/status"
	HandlerPathInternalReplicaElection = HandlerPathInternalReplica + "/leader"
	HandlerPathInternalReplicaSync     = HandlerPathInternalReplica + "/sync"

	QueryParameterForce = "force"

	DiscoveryIdLength     = 64
	LogBookEntryStorePath = model.Path("_internal/logbook_counter")
)

var (
	instanceId = util.RandomString(DiscoveryIdLength)
)

func GetInstanceId() string {
	return instanceId
}
