package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// only needed below for sample processing

func main() {

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8080")

	// accept connection on port

	for {

		conn, _ := ln.Accept()
		log.Println("Accepting Connection....")

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {

	log.Println("Polling Connectiong starting ....")

	reader := bufio.NewReader(conn)

	for {
		b, pre, err := reader.ReadLine()

		fmt.Println(b, pre, err)

		if err != nil {
			fmt.Println(err)
			log.Println("Closing socket")
			conn.Close()
		}

		log.Printf("msg : %s", b)

		conn.Write([]byte("OK\n"))

	}

}
