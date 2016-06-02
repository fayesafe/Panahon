package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"Panahon/database"
	"Panahon/logger"
	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb/client/v2"
)

type dbClient interface {
	Query(q client.Query) (*client.Response, error)
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

// queryHandleInterval is a HTTP handler for querying an interval of time
// using database API
func queryHandleInterval() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		low, ok := vars["low"]
		if !ok || !isStringNum(low) {
			low = "0"
		}

		high, ok := vars["high"]
		if !ok || !isStringNum(high) {
			high = "2147483647"
		}

		response, err := database.QueryInterval(low, high)
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
func queryHandleLast() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset, ok := vars["last"]
		if !ok || !isStringNum(offset) {
			offset = ""
		}

		response, err := database.QueryAll(offset)
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
func queryHandleAverage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset, ok := vars["offset"]
		if !ok || !isStringNum(offset) {
			offset = "0"
		}
		interval := vars["interval"]
		unit := vars["unit"]
		col := vars["col"]

		response, err := database.QueryAverage(col, interval + unit, offset)
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

// sendPayload is sending received data from query as JSON
// taking response and respective writer as input
func sendPayload(queryResponse *client.Response, w http.ResponseWriter) {
	for i := range queryResponse.Results {
		payload, err := json.Marshal(queryResponse.Results[i])
		if err != nil {
			logger.Error.Println(err)
		}
		logger.Info.Printf("Sending Payload: %s", string(payload))
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
func StartServer(appPort string, staticPath string) {
	router := mux.NewRouter()
	router.StrictSlash(false)

	logger.Info.Println("Adding routes to router/subrouter")
	addAPIRoutes(router)
	logger.Info.Printf("Serving static content of dir: %s", staticPath)
	router.PathPrefix("/").Handler(staticServe(staticPath))

	http.Handle("/", router)
	http.ListenAndServe(":"+appPort, nil)
}
