package health

import "testing"

func TestDifferentHealthIssues(t *testing.T) {
	healthy = HEALTHY

	SetUnhealthy(NON_HEALTHY_SERVER)
	SetUnhealthy(NON_HEALTHY_LEADER)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(NON_HEALTHY_SERVER)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(NON_HEALTHY_LEADER)
	if !IsHealth() {
		t.Fatal("Should be healthy")
	}
}
