package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"Panahon/database"
	"Panahon/logger"
	"github.com/influxdata/influxdb/client/v2"
)

type MockInfluxDbError struct{}
type MockInfluxDbHappy struct{}

func (db MockInfluxDbError) Query(q client.Query) (*client.Response, error) {
	return nil, errors.New("=== Test Error ===")
}

func (db MockInfluxDbHappy) Query(q client.Query) (*client.Response, error) {
	returnVal := new(client.Response)
	return returnVal, nil
}

func TestQueryOk(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleLast()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying all ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	if w.Body.String() != "" {
		t.Errorf("Incorrect Body.")
	}
}

func TestQueryError500(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleLast()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying all ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestQueryIntervalOk(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleInterval()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying range ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	if w.Body.String() != "" {
		t.Errorf("Incorrect Body.")
	}
}

func TestQueryIntervalError500(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleInterval()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying range ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestQueryAverageOk(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleAverage()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying range ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
	if w.Body.String() != "" {
		t.Errorf("Incorrect Body.")
	}
}

func TestQueryAverageError500(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	database.Init(mockInflux, "", "")

	queryHandle := queryHandleAverage()
	req, _ := http.NewRequest("GET", "", nil)
	w := httptest.NewRecorder()

	t.Log("Querying range ...")
	queryHandle.ServeHTTP(w, req)
	t.Log(w.Code)

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
