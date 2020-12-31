package context

type key int

const (
	RequestIDKey key = 0

	HandlerPathApi = "/api"

	HandlerPathStore = HandlerPathApi + "/store/"

	HandlerPathInternal                = HandlerPathApi + "/internal"
	HandlerPathInternalId              = HandlerPathInternal + "/id"
	HandlerPathInternalReplicaElection = HandlerPathInternal + "/replica/election"

	DiscoveryIdLength = 64
)
