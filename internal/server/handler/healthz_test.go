package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandlerNonHealthy(t *testing.T) {
	req := createHealthyRequest(t)
	healthy = HEALTHY

	rr := httptest.NewRecorder()
	SetUnhealthy(NON_HEALTHY_SERVER)

	Healthz().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}
}

func TestHealthCheckHandlerHealthy(t *testing.T) {
	req := createHealthyRequest(t)
	healthy = HEALTHY

	rr := httptest.NewRecorder()
	SetHealthy(HEALTHY)

	Healthz().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func createHealthyRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func TestDifferentHealthIssues(t *testing.T) {
	healthy = HEALTHY

	SetUnhealthy(NON_HEALTHY_SERVER)
	SetUnhealthy(NON_HEALTHY_ELECTION)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(NON_HEALTHY_SERVER)
	if IsHealth() {
		t.Fatal("Should be unhealthy")
	}
	SetHealthy(NON_HEALTHY_ELECTION)
	if !IsHealth() {
		t.Fatal("Should be healthy")
	}
}
