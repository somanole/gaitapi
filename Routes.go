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
        "CreateUser",
        "POST",
        "/users",
        CreateUser,
    },
	Route{
        "LoginUser",
        "POST",
        "/users/login",
        LoginUser,
    },
	Route{
        "GetUser",
        "GET",
        "/users/{userid}",
        GetUser,
    },
	Route{
        "UpdateUser",
        "PUT",
        "/users/{userid}",
        UpdateUser,
    },
	Route{
        "GetUserByEmail",
        "GET",
        "/users/email/{email}",
        GetUserByEmail,
    },
	Route{
        "GetUserExtraInfo",
        "GET",
        "/users/{userid}/extrainfo",
        GetUserExtraInfo,
    },
	Route{
        "CreateMatch",
        "POST",
        "/matches",
        CreateMatch,
    },
	Route{
        "GetUserMatch",
        "GET",
        "/users/{userid}/match",
        GetUserMatch,
    },
	Route{
        "GetUserChats",
        "GET",
        "/users/{userid}/chats",
        GetUserChats,
    },
	Route{
        "UpdateChat",
        "PUT",
        "/users/{userid}/chats/{receiverid}",
        UpdateChat,
    },
	Route{
        "CreateMessage",
        "POST",
        "/users/{userid}/messages/{receiverid}",
        CreateMessage,
    },
	Route{
        "GetUserMessagesByReceiverId",
        "GET",
        "/users/{userid}/messages/{receiverid}",
        GetUserMessagesByReceiverId,
    },
	Route{
        "CreateUserActivity",
        "POST",
        "/users/{userid}/activities",
        CreateUserActivity,
    },
	Route{
        "GetUserActivity",
        "GET",
        "/users/{userid}/activities",
        GetUserActivity,
    },
	Route{
        "CreateUserAcceleration",
        "POST",
        "/users/{userid}/accelerations",
        CreateUserAcceleration,
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
        "HelpPagePOSTUser",
        "GET",
        "/help/POST-user",
        HelpPagePOSTUser,
    },
	Route{
        "HelpPagePUTUser",
        "GET",
        "/help/PUT-user",
        HelpPagePUTUser,
    },
	Route{
        "HelpPageGETUserId",
        "GET",
        "/help/GET-user-id",
        HelpPageGETUserId,
    },
	Route{
        "HelpPageGETUserEmail",
        "GET",
        "/help/GET-user-email",
        HelpPageGETUserEmail,
    },
	Route{
        "HelpPageGETExtraInfo",
        "GET",
        "/help/GET-extrainfo",
        HelpPageGETExtraInfo,
    },
	Route{
        "HelpPageGETMatch",
        "GET",
        "/help/GET-match",
        HelpPageGETMatch,
    },
	Route{
        "HelpPagePOSTMessage",
        "GET",
        "/help/POST-message",
        HelpPagePOSTMessage,
    },
	Route{
        "HelpPageGETMessage",
        "GET",
        "/help/GET-message",
        HelpPageGETMessage,
    },
	Route{
        "HelpPagePOSTActivity",
        "GET",
        "/help/POST-activity",
        HelpPagePOSTActivity,
    },
	Route{
        "HelpPageGETActivity",
        "GET",
        "/help/GET-activity",
        HelpPageGETActivity,
    },
	Route{
        "HelpPagePOSTUserLogin",
        "GET",
        "/help/POST-user-login",
        HelpPagePOSTUserLogin,
    },
	Route{
        "HelpPagePOSTMatch",
        "GET",
        "/help/POST-match",
        HelpPagePOSTMatch,
    },
	Route{
        "HelpPageGETChats",
        "GET",
        "/help/GET-chats",
        HelpPageGETChats,
    },
	Route{
        "HelpPagePUTChats",
        "GET",
        "/help/PUT-chats",
        HelpPagePUTChats,
    },
}
