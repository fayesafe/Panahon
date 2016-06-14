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

// defineRoutes contains all routes of the API
func defineRoutes(influxClient dbClient) Routes {
	var routes = Routes{
		Route{
			"MaxVal",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleMax(influxClient),
		},
		Route{
			"MaxVal",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
			queryHandleMax(influxClient),
		},
		Route{
			"MaxVal",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
			queryHandleMax(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleAverage(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
			queryHandleAverage(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
			queryHandleAverage(influxClient),
		},
		Route{
			"RangeHighLow",
			"GET",
			"/range/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleInterval(influxClient),
		},
		Route{
			"RangeLowOnly",
			"GET",
			"/range/{low:[0-9]+}",
			queryHandleInterval(influxClient),
		},
		Route{
			"Range",
			"GET",
			"/range",
			queryHandleInterval(influxClient),
		},
		Route{
			"LastN",
			"GET",
			"/last/{last:[0-9]+}",
			queryHandleLast(influxClient),
		},
		Route{
			"Last",
			"GET",
			"/last",
			queryHandleLast(influxClient),
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
	return routes
}

// AddAPIRoutes adds routes and subroutes on router
func addAPIRoutes(router *mux.Router, routes Routes) {
	subRouter := router.PathPrefix("/api").Subrouter()
	for _, i := range routes {
		subRouter.
			Methods(i.Method).
			PathPrefix(i.Prefix).
			Name(i.Name).
			Handler(i.Handler)
	}
}
