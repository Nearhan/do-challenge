package main

import (
	"errors"
	"log"
	"os"
)

// ErrBadMsg is an error where QUERY || INDEX || ERROR is not
// in the first section of the command and it checks the Format
// as CMD|package|dep1,dep2
var ErrBadMsg = errors.New("Invalid Message Format")

// ErrPkgNotFound  returns when a Pkg doesn't exist in the store
var ErrPkgNotFound = errors.New("Pkg is not found")

func main() {

	port, ok := os.LookupEnv("TCP_PORT")

	if !ok {
		port = "8080"
	}

	server, err := MakeNewServer(port)

	if err != nil {
		log.Fatal("server unable to start, check server port")
	}

	server.Start()

}
