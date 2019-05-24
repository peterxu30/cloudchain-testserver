package main

import (
	"github.com/gorilla/mux"
)

type TestServer struct {
	router *mux.Router
}

func NewTestServer() *TestServer {
	server := &TestServer{
		router: mux.NewRouter().StrictSlash(true),
	}

	InitializeRoutes(server)

	return server
}
