package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Panahon/database"
	"Panahon/logger"
	"Panahon/station"

	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb/client/v2"
)

type dbConnection interface {
	QuerySensorData(last string) (*client.Response, error)
	QueryData(series string, startDate string, endDate string) (*client.Response, error)
	QueryAll(offset string) (*client.Response, error)
	QueryInterval(low string, high string) (*client.Response, error)
	QueryAverage(
		col string,
		interval string,
		offset string,
		end string) (*client.Response, error)
	QueryMax(
		col string,
		interval string,
		offset string,
		end string) (*client.Response, error)
	QueryMin(
		col string,
		interval string,
		offset string,
		end string) (*client.Response, error)
}

// staticServe serving front end application of the weather station
// Takes path as argument, pointing to root dir of the app
func staticServe(path string) http.Handler {
	logger.Info.Printf("Serving static content: %s on route /", path)
	return http.FileServer(http.Dir(path))
}

// apiHandler is th first serving point of the API
func apiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API is alive.")

		vars := mux.Vars(r)
		key := vars["key"]

		logger.Info.Printf("Calling route /api/%s", key)

		if len(key) > 0 {
			fmt.Fprintf(w, "\nKey: %s", key)
		}
	})
}

func queryHandleData(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		series, _ := vars["series"]
		startDate, _ := vars["startDate"]
		endDate, ok := vars["endDate"]

		startDate = "'" + startDate + "'"
		if !ok {
			endDate = startDate + " + 1d"
		} else {
			endDate = "'" + endDate + "'"
		}

		response, err := influxClient.QueryData(series, startDate, endDate)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

func queryHandleSensorData(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		last, _ := vars["last"]

		response, err := influxClient.QuerySensorData(last)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// queryHandleInterval is a HTTP handler for querying an interval of time
// using database API
func queryHandleInterval(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		low, ok := vars["low"]
		if !ok || !isStringNum(low) {
			low = "0"
		}

		high, ok := vars["high"]
		if !ok || !isStringNum(high) {
			high = "2147483647000"
		}

		response, err := influxClient.QueryInterval(low, high)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// queryHandleLast is a HTTP handler for querying all / last n entries
// using database API
func queryHandleLast(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset, ok := vars["last"]
		if !ok || !isStringNum(offset) {
			offset = ""
		}

		response, err := influxClient.QueryAll(offset)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// queryAverage returns the average for a specific interval starting from an
// offset
func queryHandleAverage(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		low, ok := vars["low"]
		if !ok || !isStringNum(low) {
			low = "0"
		}
		high, ok := vars["high"]
		if !ok || !isStringNum(high) {
			high = "2147483647000"
		}
		interval := vars["interval"]
		col := vars["col"]

		response, err := influxClient.QueryAverage(col, interval, low, high)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// queryHandleMax returns the max value for a given column
func queryHandleMax(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		low, ok := vars["low"]
		if !ok || !isStringNum(low) {
			low = "0"
		}

		high, ok := vars["high"]
		if !ok || !isStringNum(high) {
			high = "2147483647000"
		}

		interval := vars["interval"]
		col := vars["col"]

		response, err := influxClient.QueryMax(col, interval, low, high)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// queryHandleMin returns the min value for a given column
func queryHandleMin(influxClient dbConnection) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		low, ok := vars["low"]
		if !ok || !isStringNum(low) {
			low = "0"
		}

		high, ok := vars["high"]
		if !ok || !isStringNum(high) {
			high = "2147483647000"
		}

		interval := vars["interval"]
		col := vars["col"]

		response, err := influxClient.QueryMin(col, interval, low, high)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
		} else {
			sendPayload(response, w)
		}
	})
}

// handleMeasurement reads the sensors and returns when finished
func handleMeasurement(influxClient dbConnection, sensors station.Sensors) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sensors.Read(influxClient.(database.DBClient))

		w.Header().Set("Content-Type", "text/json")
		w.Write([]byte("done"))
	})
}

// sendPayload is sending received data from query as JSON
// taking response and respective writer as input
func sendPayload(queryResponse *client.Response, w http.ResponseWriter) {
	for i := range queryResponse.Results {
		payload, err := json.Marshal(queryResponse.Results[i])
		if err != nil {
			logger.Error.Println(err)
		}
		//logger.Info.Printf("Sending Payload: %s", string(payload))
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}
}

// isStringNum checks if a string is a representation of a number
func isStringNum(strToCheck string) bool {
	if _, err := strconv.Atoi(strToCheck); err != nil {
		return false
	}
	return true
}

// StartServer starts server
func StartServer(influxClient dbConnection, sensors station.Sensors, appPort string, staticPath string) {
	router := mux.NewRouter()
	router.StrictSlash(false)

	logger.Info.Println("Adding routes to router/subrouter")
	addAPIRoutes(router, defineRoutes(influxClient, sensors))
	logger.Info.Printf("Serving static content of dir: %s", staticPath)
	router.PathPrefix("/").Handler(staticServe(staticPath))

	http.Handle("/", router)
	http.ListenAndServe(":"+appPort, nil)
}
