package handler

import (
	health2 "github.com/timonback/keyvaluestore/internal/server/health"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandlerNonHealthy(t *testing.T) {
	req := createHealthyRequest(t)
	resetToHealthyState()

	rr := httptest.NewRecorder()
	health2.SetUnhealthy(health2.NON_HEALTHY_SERVER)

	Healthz().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}
}

func resetToHealthyState() {
	health2.SetHealthy(health2.NON_HEALTHY_SERVER)
	health2.SetHealthy(health2.NON_HEALTHY_LEADER)
}

func TestHealthCheckHandlerHealthy(t *testing.T) {
	req := createHealthyRequest(t)
	resetToHealthyState()

	rr := httptest.NewRecorder()
	health2.SetHealthy(health2.HEALTHY)

	Healthz().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func createHealthyRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}
