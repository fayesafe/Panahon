package database

import (
	"Panahon/logger"
	"errors"
	"github.com/influxdata/influxdb/client/v2"
	"io/ioutil"
	"reflect"
	"testing"
)

type InfluxMock interface {
	Query(q client.Query) (*client.Response, error)
}

type MockInfluxDbError struct{}
type MockInfluxDbHappy struct{}

func (db MockInfluxDbError) Query(q client.Query) (*client.Response, error) {
	//returnVal := new(client.Response)
	return nil, errors.New("=== Test Error ===")
}

func (db MockInfluxDbHappy) Query(q client.Query) (*client.Response, error) {
	returnVal := new(client.Response)
	return returnVal, nil
}

func TestQueryAllHappy(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	SetupDatabaseMock(mockInflux)
	response, err := QueryAll("42")
	t.Log("Checking Error ...")
	if err != nil {
		t.Errorf("Error thrown")
	}
	t.Log("Checking Response ...")
	if !CheckResponseEquality(response, new(client.Response)) {
		t.Errorf("Result not Equal")
	}
}

func TestQueryAllError(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	SetupDatabaseMock(mockInflux)
	response, err := QueryAll("42")
	t.Log("Checking Error ...")
	if err == nil {
		t.Errorf("No Error thrown")
	}
	t.Log("Checking Response ...")
	if CheckResponseEquality(response, new(client.Response)) {
		t.Errorf("No Empty Response")
	}
}

func TestQueryRangeHappy(t *testing.T) {
	mockInflux := MockInfluxDbHappy{}
	SetupDatabaseMock(mockInflux)
	response, err := QueryInterval("42", "42")
	t.Log("Checking Error ...")
	if err != nil {
		t.Errorf("Error thrown")
	}
	t.Log("Checking Response ...")
	if !CheckResponseEquality(response, new(client.Response)) {
		t.Errorf("Result not Equal")
	}
}

func TestQueryRangeError(t *testing.T) {
	mockInflux := MockInfluxDbError{}
	SetupDatabaseMock(mockInflux)
	response, err := QueryInterval("42", "42")
	t.Log("Checking Error ...")
	if err == nil {
		t.Errorf("No Error thrown")
	}
	t.Log("Checking Response ...")
	if CheckResponseEquality(response, new(client.Response)) {
		t.Errorf("No Empty Response")
	}
}

func TestMain(m *testing.M) {
	logger.Init(ioutil.Discard, ioutil.Discard, ioutil.Discard)
	m.Run()
}

func CheckResponseEquality(
	respObj *client.Response,
	compObj *client.Response) bool {
	if respObj == nil {
		return false
	}
	if !(reflect.DeepEqual(respObj.Results, compObj.Results)) ||
		respObj.Err != compObj.Err {
		return false
	}
	return true
}

func SetupDatabaseMock(databaseMock InfluxMock) {
	Init(databaseMock)
}
