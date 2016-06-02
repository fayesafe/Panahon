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

// StaticServe serving front end application of the weather station
// Takes path as argument, pointing to root dir of the app
func StaticServe(path string) http.Handler {
	logger.Info.Printf("Serving static content: %s on route /", path)
	return http.FileServer(http.Dir(path))
}

// APIHandler is th first serving point of the API
func APIHandler() http.Handler {
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

// QueryHandleInterval is a HTTP handler for querying an interval of time
// using database API
func QueryHandleInterval() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		low, ok := vars["low"]
		if !ok || !IsStringNum(low) {
			low = "0"
		}

		high, ok := vars["high"]
		if !ok || !IsStringNum(high) {
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
			SendPayload(response, w)
		}
	})
}

// QueryHandle is a HTTP handler for querying all / last n entries
// using database API
func QueryHandleLast() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset := vars["last"]
		if !IsStringNum(offset) {
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
			SendPayload(response, w)
		}
	})
}

// SendPayload is sending received data from query as JSON
// taking response and respective writer as input
func SendPayload(queryResponse *client.Response, w http.ResponseWriter) {
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

// IsStringNum checks if a string is a representation of a number
func IsStringNum(strToCheck string) bool {
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
	AddAPIRoutes(router)
	logger.Info.Printf("Serving static content of dir: %s", staticPath)
	router.PathPrefix("/").Handler(StaticServe(staticPath))

	http.Handle("/", router)
	http.ListenAndServe(":"+appPort, nil)
}
