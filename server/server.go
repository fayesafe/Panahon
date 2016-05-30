package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"panahon/logger"
)

type Env struct {
	Client client.Client
}

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

func QueryHandle(influxClient dbClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var q client.Query

		if index, ok := mux.Vars(r)["last"]; ok {
			q = client.NewQuery(
				"SELECT * FROM meas ORDER BY DESC LIMIT "+index,
				"test",
				"s")
			logger.Info.Println("Getting last " + index + " entries")
		} else {
			q = client.NewQuery("SELECT * FROM meas", "test", "s")
			logger.Info.Println("Calling route /api/get")
		}

		response, err := influxClient.Query(q)
		if err != nil {
			logger.Error.Println(err)
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
			return
		} else if response.Error() != nil {
			logger.Error.Println(response.Error())
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError)
			return
		} else {
			for i := range response.Results {
				payload, err := json.Marshal(response.Results[i])
				if err != nil {
					logger.Error.Println(err)
				}

				logger.Info.Println("Sending Payload: " + string(payload))
				w.Header().Set("Content-Type", "application/json")
				w.Write(payload)
			}
		}
	})
}

func AddApiRoutes(influxClient dbClient, router *mux.Router) {
	logger.Info.Println("Adding routes to router/subrouter")

	// Subroutes on /api go here
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.PathPrefix("/get/{last:[0-9]+}").Methods("GET").Handler(QueryHandle(influxClient))
    subRouter.PathPrefix("/get").Methods("GET").Handler(QueryHandle(influxClient))
	subRouter.PathPrefix("/{key}").Methods("GET").Handler(ApiHandler())
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
	logger.Info.Println("InfluxDB client initialized")

    router := mux.NewRouter()
	router.StrictSlash(false)

	AddApiRoutes(influxClient, router)

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}
