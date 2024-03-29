package main

import (
	"fmt"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// InitializeRoutes initializes all routes for the given TestServer.
func InitializeRoutes(s *TestServer) {
	var routes = Routes{
		Route{
			"Index",
			"GET",
			"/",
			handleDefault(),
		},
		Route{
			"AddBlocksTest",
			"POST",
			"/addblockstest",
			s.AddFiftyBlocksTest(),
		},
		Route{
			"AddAndReadBlocksTest",
			"POST",
			"/addandreadblockstest",
			s.SimultaneouslyAddAndReadFiftyBlocksTest(),
		},
		Route{
			"Reset",
			"POST",
			"/reset",
			s.Reset(),
		},
	}

	for _, route := range routes {
		s.router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)

	}
}

func handleDefault() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "CloudChain Test Server")
	}
}
