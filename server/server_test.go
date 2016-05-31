package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Panahon/logger"
	"github.com/influxdata/influxdb/client/v2"
)

type MockInfluxDbError struct{}
type MockInfluxDbHappy struct{}

func (db MockInfluxDbError) Query(q client.Query) (*client.Response, error) {
	returnVal := new(client.Response)
	return returnVal, errors.New("=== Test Error ===")
}

func (db MockInfluxDbHappy) Query(q client.Query) (*client.Response, error) {
	returnVal := new(client.Response)
	return returnVal, nil
}

func TestQueryOk(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	queryHandle := QueryHandle(mockInflux)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	queryHandle.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	if w.Body.String() != "" {
		t.Errorf("Incorrect Body.")
	}
}

func TestQueryError500(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	queryHandle := QueryHandle(mockInflux)
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()
	queryHandle.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestMain(m *testing.M) {
	SetupLogger()
	exitStatus := m.Run()
	os.Exit(exitStatus)
}

func SetupLogger() {
	logger.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard)
}
