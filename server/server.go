package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"Panahon/database"
	"Panahon/logger"
	"github.com/gorilla/mux"
	"github.com/influxdata/influxdb/client/v2"
)

type dbClient interface {
	Query(q client.Query) (*client.Response, error)
}

func StaticServe(path string) http.Handler {
	logger.Info.Println("Serving static content: " + path + " on route /")
	return http.FileServer(http.Dir(path))
}

func ApiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API is alive.")

		vars := mux.Vars(r)
		key := vars["key"]

		logger.Info.Println("Calling route /api/" + key)

		if len(key) > 0 {
			fmt.Fprintf(w, "\nKey: %s", key)
		}
	})
}

func QueryHandleInterval() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		low, ok := vars["low"]
		_, errLow := strconv.Atoi(low)
		if !ok || errLow != nil {
			low = "0"
		}

		high, ok := vars["high"]
		_, errHigh := strconv.Atoi(high)
		if !ok || errHigh != nil {
			high = "2147483647"
		}

		response, err := database.QueryInterval(low, high)
		SendPayload(response, err, w)
	})
}

func QueryHandle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		offset := vars["last"]
		if _, errOffset := strconv.Atoi(offset); errOffset != nil {
			offset = ""
		}
		response, err := database.QueryAll(offset)
		SendPayload(response, err, w)
	})
}

func SendPayload(queryResponse *client.Response, err error, w http.ResponseWriter) {
	if err != nil {
		logger.Error.Println(err)
		http.Error(
			w,
			"Internal Server Error",
			http.StatusInternalServerError)
		return
	} else {
		for i := range queryResponse.Results {
			payload, err := json.Marshal(queryResponse.Results[i])
			if err != nil {
				logger.Error.Println(err)
			}
			logger.Info.Println("Sending Payload: " + string(payload))
			w.Header().Set("Content-Type", "application/json")
			w.Write(payload)
		}
	}
}

func AddApiRoutes(router *mux.Router) {
	logger.Info.Println("Adding routes to router/subrouter")

	// Subroutes on /api go here
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.PathPrefix(
		"/range/{low:[0-9]+}/{high:[0-9]+}").Methods("GET").Handler(
		QueryHandleInterval())
	subRouter.PathPrefix(
		"/range/{low:[0-9]+}").Methods("GET").Handler(
		QueryHandleInterval())
	subRouter.PathPrefix(
		"/range/").Methods("GET").Handler(
		QueryHandleInterval())
	subRouter.PathPrefix(
		"/get/{last:[0-9]+}").Methods("GET").Handler(QueryHandle())
	subRouter.PathPrefix(
		"/get").Methods("GET").Handler(QueryHandle())
	subRouter.PathPrefix(
		"/{key}").Methods("GET").Handler(ApiHandler())
	subRouter.Methods("GET").Handler(ApiHandler())

	// Routes on Router go here
	router.PathPrefix("/").Handler(StaticServe("./app/"))
}

func StartServer() {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		logger.Error.Println(err)
		os.Exit(1)
	}
	database.Init(influxClient)
	logger.Info.Println("InfluxDB client initialized")

	router := mux.NewRouter()
	router.StrictSlash(false)

	AddApiRoutes(router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
