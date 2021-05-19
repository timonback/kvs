package health

import "testing"

func TestDifferentHealthIssues(t *testing.T) {
	healthy = HEALTHY

	SetUnhealthy(SERVER_STATUS)
	SetUnhealthy(REPLICA_STATUS)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(SERVER_STATUS)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(REPLICA_STATUS)
	if !IsHealth() {
		t.Fatal("Should be healthy")
	}
}
