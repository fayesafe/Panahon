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
		"Average",
		"GET",
		"/av/{col:[a-z]+}/{interval:[0-9]+}/{unit:(ms)|[usmhdw]}/{offset:[0-9]+}",
		queryHandleAverage(),
	},
	Route{
		"Average",
		"GET",
		"/av/{col:[a-z]+}/{interval:[0-9]+}/{unit:(ms)|[usmhdw]}",
		queryHandleAverage(),
	},
	Route{
		"RangeHighLow",
		"GET",
		"/range/{low:[0-9]+}/{high:[0-9]+}",
		queryHandleInterval(),
	},
	Route{
		"RangeLowOnly",
		"GET",
		"/range/{low:[0-9]+}",
		queryHandleInterval(),
	},
	Route{
		"Range",
		"GET",
		"/range",
		queryHandleInterval(),
	},
	Route{
		"LastN",
		"GET",
		"/last/{last:[0-9]+}",
		queryHandleLast(),
	},
	Route{
		"Last",
		"GET",
		"/last",
		queryHandleLast(),
	},
	Route{
		"GeneralApiWithKey",
		"GET",
		"/{key}",
		apiHandler(),
	},
	Route{
		"GeneralApi",
		"GET",
		"/",
		apiHandler(),
	},
}

// AddAPIRoutes adds routes and subroutes on router
func addAPIRoutes(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	for _, i := range routes {
		subRouter.
			Methods(i.Method).
			PathPrefix(i.Prefix).
			Name(i.Name).
			Handler(i.Handler)
	}
}
