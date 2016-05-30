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

func StaticServe(path string) http.Handler {
	logger.Info.Println("Serving static content: " + path + " on route /")
	return http.FileServer(http.Dir(path))
}

func ApiHandler(env *Env) http.Handler {
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

func QueryHandle(env *Env) http.Handler {
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

		response, err := env.Client.Query(q)
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

func StartServer() {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		logger.Error.Println(err)
		os.Exit(1)
	}
	logger.Info.Println("InfluxDB client initialized")

	env := &Env{Client: c}

	r := mux.NewRouter()
	r.StrictSlash(false)

	s := r.PathPrefix("/api").Subrouter()

	s.PathPrefix("/get/{last:[0-9]+}").Methods("GET").Handler(QueryHandle(env))
	s.PathPrefix("/get").Methods("GET").Handler(QueryHandle(env))

	s.PathPrefix("/{key}").Methods("GET").Handler(ApiHandler(env))
	s.Methods("GET").Handler(ApiHandler(env))

	r.PathPrefix("/").Handler(StaticServe("./app/"))

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
