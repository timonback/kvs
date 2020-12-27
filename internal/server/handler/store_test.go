package handler

import (
	"encoding/json"
	"github.com/timonback/keyvaluestore/internal/server/context"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	storeTestImpl store2.Service = store2.NewStoreInmemoryService()
)

func newRequest(t *testing.T, method string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, context.HandlerPathStore+"key", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newGetRequest(t *testing.T) *http.Request {
	return newRequest(t, "GET", nil)
}

func newPostRequest(t *testing.T, request storeRequestPost) *http.Request {
	requestJson, _ := json.Marshal(request)
	return newRequest(t, "POST", strings.NewReader(string(requestJson)))
}

func TestStoreHandlerGetNonExisting(t *testing.T) {
	rr := httptest.NewRecorder()

	Store(storeTestImpl).ServeHTTP(rr, newGetRequest(t))
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
func TestStoreHandlerGet(t *testing.T) {
	rr := httptest.NewRecorder()

	requestPost := storeRequestPost{
		Content: "CONTENT",
	}
	Store(storeTestImpl).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	Store(storeTestImpl).ServeHTTP(rr2, newGetRequest(t))
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	getResponse := storeResponseGet{
		Key:     "key",
		Content: requestPost.Content,
	}
	if body, _ := json.Marshal(getResponse); rr2.Body.String() != string(body) {
		t.Errorf("response did not match: got %v want %v", rr2.Body.String(), string(body))
	}
}
