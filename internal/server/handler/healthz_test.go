package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	SetHealthy(NON_HEALTHY)

	Healthz().ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusServiceUnavailable {
		t.Errorf("hanlder reqturned wrong status code: got %v want %v", status, http.StatusServiceUnavailable)
	}

	rr2 := httptest.NewRecorder()
	SetHealthy(HEALTHY)

	Healthz().ServeHTTP(rr2, req)
	if status := rr2.Code; status != http.StatusNoContent {
		t.Errorf("hanlder reqturned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
