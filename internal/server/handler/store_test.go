package handler

import (
	"encoding/json"
	"github.com/timonback/keyvaluestore/internal/server/context"
	"github.com/timonback/keyvaluestore/internal/server/handler/pojo"
	store2 "github.com/timonback/keyvaluestore/internal/store"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newRequest(t *testing.T, method string, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, context.HandlerPathStore+path, body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func newListRequest(t *testing.T) *http.Request {
	return newRequest(t, "GET", "", nil)
}

func newGetRequest(t *testing.T) *http.Request {
	return newRequest(t, "GET", "key", nil)
}

func newPostRequest(t *testing.T, request pojo.StoreRequestPost) *http.Request {
	requestJson, _ := json.Marshal(request)
	return newRequest(t, "POST", "key", strings.NewReader(string(requestJson)))
}

func newPutRequest(t *testing.T, request pojo.StoreRequestPost) *http.Request {
	requestJson, _ := json.Marshal(request)
	return newRequest(t, "PUT", "key", strings.NewReader(string(requestJson)))
}

func newDeleteRequest(t *testing.T) *http.Request {
	return newRequest(t, "DELETE", "key", nil)
}

func TestStoreHandlerList(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	Store(store).ServeHTTP(rr, newListRequest(t))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
func TestStoreHandlerGetNonExisting(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	Store(store).ServeHTTP(rr, newGetRequest(t))
	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestStoreHandlerGet(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "CONTENT",
	}
	Store(store).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	Store(store).ServeHTTP(rr2, newGetRequest(t))
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

func TestStoreHandlerDeleteNonExisting(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	Store(store).ServeHTTP(rr, newDeleteRequest(t))
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestStoreHandlerDelete(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "Content",
	}
	Store(store).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	Store(store).ServeHTTP(rr2, newDeleteRequest(t))
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestStoreHandlerPost(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "Content",
	}
	Store(store).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}
}

func TestStoreHandlerPostOnExistingElement(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "Content",
	}
	Store(store).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	Store(store).ServeHTTP(rr2, newPostRequest(t, requestPost))
	if status := rr2.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusBadRequest, rr2.Body.String())
	}
}

func TestStoreHandlerPutOnNonExistingElement(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "Content",
	}
	Store(store).ServeHTTP(rr, newPutRequest(t, requestPost))
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusBadRequest, rr.Body.String())
	}
}

func TestStoreHandlerPut(t *testing.T) {
	rr := httptest.NewRecorder()
	store := store2.NewStoreInmemoryService("")

	requestPost := pojo.StoreRequestPost{
		Content: "Content",
	}
	Store(store).ServeHTTP(rr, newPostRequest(t, requestPost))
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr.Body.String())
	}

	rr2 := httptest.NewRecorder()
	Store(store).ServeHTTP(rr2, newPutRequest(t, requestPost))
	if status := rr2.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v (message %v)", status, http.StatusOK, rr2.Body.String())
	}
}
