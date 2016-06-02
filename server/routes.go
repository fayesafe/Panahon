package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name    string
	Method  string
	Prefix  string
	Handler http.Handler
}

type Routes []Route

// routes contains all routes of the API
var routes = Routes{
	Route{
		"RangeHighLow",
		"GET",
		"/range/{low:[0-9]+}/{high:[0-9]+}",
		QueryHandleInterval(),
	},
	Route{
		"RangeLowOnly",
		"GET",
		"/range/{low:[0-9]+}",
		QueryHandleInterval(),
	},
	Route{
		"Range",
		"GET",
		"/range",
		QueryHandleInterval(),
	},
	Route{
		"LastN",
		"GET",
		"/last/{last:[0-9]+}",
		QueryHandleLast(),
	},
	Route{
		"Last",
		"GET",
		"/last",
		QueryHandleLast(),
	},
	Route{
		"GeneralApiWithKey",
		"GET",
		"/{key}",
		APIHandler(),
	},
	Route{
		"GeneralApi",
		"GET",
		"/",
		APIHandler(),
	},
}

// AddAPIRoutes adds routes and subroutes on router
func AddAPIRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	for _, i := range routes {
		subRouter.
			Methods(i.Method).
			PathPrefix(i.Prefix).
			Name(i.Name).
			Handler(i.Handler)
	}
}
