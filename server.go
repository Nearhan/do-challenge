package main

import (
	"bufio"
	"log"
	"net"
	"strings"
	"sync"
)

// Msg ...
type Msg struct {
	Command string
	Package string
	Deps    []string
}

// Server struct for server implementation
type Server struct {
	PkgStore *PkgStore

	Listener net.Listener
}

// MakeNewServer creates a new Server Struct and parses port
func MakeNewServer(port string) (*Server, error) {

	pkgStore := &PkgStore{&sync.RWMutex{}, make(PkgDtl)}

	addr := ":" + port

	ln, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err

	}

	log.Printf("Server up .... Listening on port %s", port)

	return &Server{pkgStore, ln}, nil

}

// Start tcp server begins accepting connections
func (s *Server) Start() error {

	for {

		conn, err := s.Listener.Accept()
		log.Println("Accepting Socket Connection...")
		if err != nil {
			log.Println("server unable to accept connection")
			return err
		}

		go s.handleConnection(conn)

	}

}

func (s *Server) handleConnection(conn net.Conn) {

	reader := bufio.NewReader(conn)

	for {

		buf, _, err := reader.ReadLine()

		// check to see if EOF is sent
		if err != nil {
			log.Println("Recieved EOF on socket, closing connection")
			conn.Close()
			return
		}

		// cast to string
		raw := string(buf[:])

		// parse message
		msg, err := parseMessage(raw)

		// if err send Error back to the client
		if err == ErrBadMsg {
			sendError(conn)

		} else {

			// router for requests
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

// Sever Methods

// handleRemove deals with handling a remove request
func (s *Server) handleRemove(msg *Msg, conn net.Conn) {

	ok := s.PkgStore.Remove(msg.Package)
	if !ok {
		sendFail(conn)
		return
	}
	sendOk(conn)
}

// handleIndex deals with handling an index request
func (s *Server) handleIndex(msg *Msg, conn net.Conn) {

	ok := s.PkgStore.Add(msg)
	if !ok {
		sendFail(conn)
		return
	}
	sendOk(conn)

}

// handleQuery deals with handling query request
func (s *Server) handleQuery(msg *Msg, conn net.Conn) {

	ok := s.PkgStore.Get(msg.Package)
	if !ok {
		sendFail(conn)
		return
	}
	sendOk(conn)

}

// helper methods

// sends OK on client
func sendOk(conn net.Conn) {
	conn.Write([]byte("OK\n"))

}

// sends FAIL on client
func sendFail(conn net.Conn) {
	conn.Write([]byte("FAIL\n"))
}

// sends ERROR on client
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
	msg := &Msg{cmd, strings.ToLower(c[1]), make([]string, 0)}

	if containsDeps(c) {
		msg.Deps = strings.Split(strings.ToLower(c[2]), ",")

	}

	return msg, nil

}

// validate incoming CMD
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
