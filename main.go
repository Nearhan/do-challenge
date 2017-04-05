package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// 	r, _ := regexp.Compile(`^(INDEX|REMOVE|QUERY)\|([\w\d\-_\+]+)\|([\w\d,\-_\+]+)*\n`)

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

	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGHUP, syscall.SIGQUIT)
	exit := longRunningGoroutine(quit, server)
	log.Println("Waiting for signal")
	code := <-exit
	log.Printf("Exited with code %v\n", code)
	os.Exit(code)

}

func longRunningGoroutine(quit chan os.Signal, s *Server) chan int {
	exit := make(chan int)

	go func() {

		for {
			select {
			case <-quit:
				writeFile(s)
				exit <- 1
				break
			}
		}
	}()

	return exit
}

func writeFile(s *Server) {

	data, err := json.Marshal(s.PkgStore.Index)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile("data.text", data, 0755)
	if err != nil {
		log.Fatal(err)
	}

}
