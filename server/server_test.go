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
	"github.com/influxdata/influxdb/models"
	"github.com/drewolson/testflight"
)

type MockInfluxDbError struct{}
type MockInfluxDbHappy struct{}

var expectedResponse string = `{"Series":[{"name":"Test","tags":{"test":"test"},"columns":["test","test"]}],"Messages":null,"error":"test"}`

func (db MockInfluxDbError) Query(q client.Query) (*client.Response, error) {
	return nil, errors.New("=== Test Error ===")
}

func (db MockInfluxDbHappy) Query(q client.Query) (*client.Response, error) {
	tags := map[string]string{
		"test":"test",
	}
	result := models.Row{
		Name: "Test",
		Tags: tags,
		Columns: []string{"test", "test"},
		Values: nil,
		Err: nil,
	}
	results := []client.Result{
		client.Result{
			Series: []models.Row{result},
			Messages: nil,
			Err: "test",
		},
	}
	returnVal := client.Response{
		Results: results,
		Err: "no error",
	}
	return &returnVal, nil
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
	if w.Body.String() != expectedResponse {
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
	if w.Body.String() != expectedResponse {
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
	if w.Body.String() != expectedResponse {
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

func TestQueryLastHTTPOk(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	database.Init(mockInflux, "", "")

	testflight.WithServer(queryHandleLast(), func(r *testflight.Requester) {
        response := r.Get("/last")
		t.Log("Checking Response")
		if response.StatusCode != 200 {
			t.Errorf("Status Code not 200")
		}
		if response.Body != expectedResponse {
			t.Errorf("Wrong Body")
		}
    })
}

func TestQueryLastHTTPError(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	database.Init(mockInflux, "", "")

	testflight.WithServer(queryHandleLast(), func(r *testflight.Requester) {
        response := r.Get("/last")
		t.Log("Checking Response")
		if response.StatusCode != 500 {
			t.Errorf("Status Code not 500")
		}
		if response.Body != "Internal Server Error\n" {
			t.Errorf("Body not empty")
		}
    })
}

func TestMain(m *testing.M) {
	SetupLogger()
	exitStatus := m.Run()
	os.Exit(exitStatus)
}

func SetupLogger() {
	logger.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard)
}
