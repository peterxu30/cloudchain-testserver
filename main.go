package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"google.golang.org/appengine"
)

func init() {
	log.Print("Prestigechain server started.")
	testServerInit()
}

var _s *TestServer

func testServerInit() {
	_s = NewTestServer()

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	configuredRouter := handlers.CORS(allowedOrigins, allowedMethods)(_s.router)
	http.Handle("/", configuredRouter)

	// Initialize the testCloudChain singleton.
	GetTestCloudChain(context.Background())
}

func main() {
	appengine.Main()
}
