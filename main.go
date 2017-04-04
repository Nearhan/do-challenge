package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// 	r, _ := regexp.Compile(`^(INDEX|REMOVE|QUERY)\|([\w\d\-_\+]+)\|([\w\d,\-_\+]+)*\n`)

// ErrBadMsg is an error where QUERY || INDEX || ERROR is not
// in the first section of the command and it checks the Format
// as CMD|package|dep1,dep2
var ErrBadMsg = errors.New("Invalid Message Format")

// ErrPkgNotFound  returns when a Pkg doesn't exist in the store
var ErrPkgNotFound = errors.New("Pkg is not found")

// Msg ...
type Msg struct {
	Command string
	Package string
	Deps    []string
}

// PkgNm simple type alias for string, the name of the package
type PkgNm string

// PkgDetail ...
type PkgDtl struct {

	// Package dependencies
	Deps []string

	// What package require this one
	Reqs []string
}

// PkgStore is the representation of that state of the package
type PkgStore struct {

	// mutex for locking
	mutex *sync.Mutex

	// list of that state
	Index map[string]PkgDtl
}

// Remove ...
func (pkSt *PkgStore) Remove(pkgName string) {

	pkSt.mutex.Lock()
	delete(pkSt.Index, pkgName)
	pkSt.mutex.Unlock()

}

// Config is the configuration for the server
type Config struct {
	port  string
	debug bool
}

type Server struct {
	PkgStore *PkgStore

	Listener net.Listener

	Counter *Counter
}

func MakeNewServer(port string) (*Server, error) {

	pkgStore := &PkgStore{&sync.Mutex{}, make(map[string]PkgDtl)}

	addr := ":" + port

	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err

	}

	log.Printf("Server up .... Listening on port %s", port)

	return &Server{pkgStore, ln, &Counter{0, &sync.Mutex{}}}, nil

}

func (s *Server) Start() error {

	for {

		conn, err := s.Listener.Accept()
		log.Println("Accepting Socket Connection...")
		if err != nil {
			log.Fatal("server unable to accept connection")
		}

		go s.handleConnection(conn)
		//go s.test(conn)

	}

}

func (s *Server) test(conn net.Conn) {

	for {
		buf := make([]byte, 1024)

		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)
		if err != nil {
			conn.Close()

		}

		sendOk(conn)

	}

}

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

func (s *Server) handleConnection(conn net.Conn) {

	for {

		// make buffer
		buf := make([]byte, 1024)

		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)

		// check to see if EOF is sent
		if err != nil {
			fmt.Println(err)
			log.Println("Recieved EOF on socket, closing connection")
			conn.Close()
			return
		}

		// cast to string
		raw := string(buf[:reqLen])

		// parse message
		msg, err := parseMessage(raw)

		if err == ErrBadMsg {
			//s.Counter.Inc()
			sendOk(conn)
			//sendError(conn)

		} else {

			switch msg.Command {
			case "REMOVE":
				s.handleRemove(msg, conn)
			case "INDEX":
				s.handleIndex(msg, conn)
			case "QUERY":
				s.handleQuery(msg, conn)
			}
		}
	}

}

func sendOk(conn net.Conn) {
	conn.Write([]byte("OK\n"))

}

func sendFail(conn net.Conn) {
	conn.Write([]byte("FAIL\n"))
}

func sendError(conn net.Conn) {
	conn.Write([]byte("ERROR\n"))
}

// parse the incoming message
func parseMessage(raw string) (*Msg, error) {

	// split on delimiter
	r := strings.TrimSpace(raw)
	c := strings.Split(r, "|")

	// check for correct length
	if len(c) != 3 {
		return nil, ErrBadMsg

	}

	// check to see command is correct
	cmd := strings.ToUpper(c[0])
	if !validCmd(cmd) {
		return nil, ErrBadMsg
	}

	// construct message
	msg := &Msg{cmd, c[1], nil}

	if containsDeps(c) {
		msg.Deps = strings.Split(c[2], ",")

	}

	return msg, nil

}

// Helper methods

func validCmd(cmd string) bool {
	switch cmd {
	case "REMOVE":
		return true
	case "QUERY":
		return true
	case "INDEX":
		return true
	default:
		return false
	}

}

func containsDeps(r []string) bool {

	if len(r[2]) > 1 {
		return true
	}
	return false
}

// Bussiness Logic Controllers

func (s *Server) handleRemove(msg *Msg, conn net.Conn) {

	pkg, ok := s.PkgStore.Index[msg.Package]
	// if pkg doesn't exist return OK

	if !ok {
		sendOk(conn)
	}

	if len(pkg.Deps) == 0 && len(pkg.Reqs) == 0 {
		s.PkgStore.Remove(msg.Package)
		sendOk(conn)

	} else {
		// can't remove becase of dependencies
		sendFail(conn)
	}

}

func (s *Server) handleIndex(msg *Msg, conn net.Conn) {

}

func (s *Server) handleQuery(msg *Msg, conn net.Conn) {

}

func splitString(raw string) []string {

	return nil

}

type Counter struct {
	Count int
	Mutex *sync.Mutex
}

func (c *Counter) Inc() {
	c.Mutex.Lock()
	c.Count++
	fmt.Println(c.Count)
	c.Mutex.Unlock()
}
