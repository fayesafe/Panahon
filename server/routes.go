package server

import (
	"Panahon/database"
	"Panahon/station"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Prefix  string
	Handler http.Handler
}

type Routes []Route

// defineRoutes contains all routes of the API
func defineRoutes(influxClient database.DBClient, sensors station.Sensors) Routes {
	var routes = Routes{
		Route{
			"MaxValBothVars",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleMax(influxClient),
		},
		Route{
			"MaxValOnlyLow",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
			queryHandleMax(influxClient),
		},
		Route{
			"MaxValAll",
			"GET",
			"/max/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
			queryHandleMax(influxClient),
		},
		Route{
			"MinValBothVars",
			"GET",
			"/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleMin(influxClient),
		},
		Route{
			"MinValOnlyLow",
			"GET",
			"/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
			queryHandleMin(influxClient),
		},
		Route{
			"MinValAll",
			"GET",
			"/min/{col:[a-z]+}/{interval:[0-9]+((ms)|[usmhdw])}",
			queryHandleMin(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z\\*]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}/{high:[0-9]+}",
			queryHandleAverage(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z*]+}/{interval:[0-9]+((ms)|[usmhdw])}/{low:[0-9]+}",
			queryHandleAverage(influxClient),
		},
		Route{
			"Average",
			"GET",
			"/av/{col:[a-z*]+}/{interval:[0-9]+((ms)|[usmhdw])}",
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
			"Measurement",
			"GET",
			"/measure",
			handleMeasurement(influxClient, sensors),
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
