// Routes test
package main

import (
    "net/http"
    "github.com/gorilla/mux"
)

type Route struct {
    Name        string
    Method      string
    Pattern     string
    HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

    router := mux.NewRouter().StrictSlash(true)
    for _, route := range routes {
        var handler http.Handler

        handler = route.HandlerFunc
        handler = Logger(handler, route.Name)

        router.
            Methods(route.Method).
            Path(route.Pattern).
            Name(route.Name).
            Handler(handler)
    }

    return router
}

var routes = Routes{
    Route{
        "Index",
        "GET",
        "/",
        Index,
    },
	Route{
        "AccelerationIndex",
        "GET",
        "/accelerations",
        AccelerationIndex,
    },
	Route{
        "AccelerationsCount",
        "GET",
        "/accelerations/count",
        AccelerationsCount,
    },
	Route{
        "AccelerationCreate",
        "POST",
        "/acceleration",
        AccelerationCreate,
    },
    Route{
        "ValidateAccessCode",
        "POST",
        "/accesscode/validate",
        ValidateAccessCode,
    },
	Route{
        "GetAccessCode",
        "GET",
        "/accesscode",
        GetAccessCode,
    },
	Route{
        "HelpPageIndex",
        "GET",
        "/help",
        HelpPageIndex,
    },
	Route{
        "HelpPageCss",
        "GET",
        "/help/css",
        HelpPageCss,
    },
	Route{
        "HelpPagePOSTAccesscodeValidate",
        "GET",
        "/help/POST-accesscode-validate",
        HelpPagePOSTAccesscodeValidate,
    },
	Route{
        "HelpPageGETAccesscode",
        "GET",
        "/help/GET-accesscode",
        HelpPageGETAccesscode,
    },
	Route{
        "HelpPagePOSTAcceleration",
        "GET",
        "/help/POST-acceleration",
        HelpPagePOSTAcceleration,
    },
	Route{
        "HelpPageGETAccelerations",
        "GET",
        "/help/GET-accelerations",
        HelpPageGETAccelerations,
    },
	Route{
        "HelpPageGETAccelerationsCount",
        "GET",
        "/help/GET-accelerations-count",
        HelpPageGETAccelerationsCount,
    },
	Route{
        "CreateUser",
        "POST",
        "/user",
        CreateUser,
    },
	Route{
        "GetUser",
        "GET",
        "/user/{id}",
        GetUser,
    },
	Route{
        "UpdateUser",
        "PUT",
        "/user/{id}",
        UpdateUser,
    },
}
