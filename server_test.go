package main

import (
	"bufio"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"testing"
)

var testPort = 8081

func TestParseMessage(t *testing.T) {

	var testMsgs = []struct {
		in   string // input string
		out1 *Msg   // output Msg Pointer
		out2 error  // output error

	}{
		{"REMOVE|git|\n", &Msg{"REMOVE", "git", makeEmptySliceStr()}, nil},
		{"QUERY|git|\n", &Msg{"QUERY", "git", makeEmptySliceStr()}, nil},
		{"INDEX|leveldb|git,mysql\n", &Msg{"INDEX", "leveldb", []string{"git", "mysql"}}, nil},
		{"TEST", nil, ErrBadMsg},
		{"REMOVE|TEST|git|osx", nil, ErrBadMsg},
		{"INDEX|git osx", nil, ErrBadMsg},
		{"slakdjflsadjfljsdlafjlsjadlfjsaldfjlsf$j\n", nil, ErrBadMsg},
		{"REMOVE|some_package_awesome_pants|\n", &Msg{"REMOVE", "some_package_awesome_pants", makeEmptySliceStr()}, nil},
	}

	for _, tt := range testMsgs {
		o1, o2 := parseMessage(tt.in)
		if !reflect.DeepEqual(o1, tt.out1) {
			t.Fatalf(" Test Parse Test failure | Expected: %s, %s | Acutal: %s, %s", tt.out1, tt.out2, o1, o2)

		}
	}

}

// TestServer is the integration tests for the tcp server
func TestServer(t *testing.T) {

	var testCases = []struct {
		name  string              // name of test case
		setup func(p int) *Server // function to configure & start server
		in    []string            // inputs that client will send
		out   []string            // output the client expects
	}{
		{
			"Server Integration Test REMOVE",
			func(p int) *Server {

				port := strconv.Itoa(p)

				s, err := MakeNewServer(port)
				if err != nil {
					t.Fatalf("Test Server unable to start on port %s", 8081, err)

				}
				// setup store state
				s.PkgStore = makeStoreWithState(&PkgDtl{
					"git": []string{},
					"vim": []string{},
					"osx": []string{"vim"},
				})

				go s.Start()
				return s

			},
			[]string{
				"REMOVE|git|\n",
				"REMOVE|vim|\n",
				"REMOVE|osx|\n",
				"REMOVE|vim|\n",
			},
			[]string{
				"OK",
				"FAIL",
				"OK",
				"OK",
			},
		},
		{
			"Server Integration Test INDEX",
			func(p int) *Server {

				port := strconv.Itoa(p)

				s, err := MakeNewServer(port)
				if err != nil {
					t.Fatalf("Test Server unable to start on port %s", testPort, err)

				}
				// setup store state
				s.PkgStore = makeStore()

				go s.Start()
				return s

			},
			[]string{
				"INDEX|git|\n",
				"INDEX|vim|osx,git\n",
				"INDEX|osx|\n",
				"INDEX|vim|osx,git\n",
			},
			[]string{
				"OK",
				"FAIL",
				"OK",
				"OK",
			},
		},
		{
			"Server Integration Test QUERY",
			func(p int) *Server {

				port := strconv.Itoa(p)

				s, err := MakeNewServer(port)
				if err != nil {
					t.Fatalf("Test Server unable to start on port %s", 8081, err)

				}
				// setup store state
				s.PkgStore = makeStoreWithState(&PkgDtl{
					"git": []string{},
					"vim": []string{},
				})

				go s.Start()
				return s

			},
			[]string{
				"QUERY|git|\n",
				"QUERY|vim|\n",
				"QUERY|osx|\n",
				"QUERY|dog|\n",
			},
			[]string{
				"OK",
				"OK",
				"FAIL",
				"FAIL",
			},
		},
		{
			"Server Integration Test ALL MESSAGES",
			func(p int) *Server {

				port := strconv.Itoa(p)

				s, err := MakeNewServer(port)
				if err != nil {
					t.Fatalf("Test Server unable to start on port %s", 8081, err)

				}
				// setup store state
				s.PkgStore = makeStoreWithState(&PkgDtl{
					"git": []string{},
					"vim": []string{},
				})

				go s.Start()
				return s

			},
			[]string{
				"REMOVE|git|\n",
				"QUERY|git|\n",
				"SDFSD|DKFJLSJ|\n",
				"INDEX|git|vim,osx\n",
				"QUERY|vim|\n",
				"INDEX|osx|\n",
				"INDEX|git|vim,osx\n",
				"REMOVE|vim|\n",
			},
			[]string{
				"OK",
				"FAIL",
				"ERROR",
				"FAIL",
				"OK",
				"OK",
				"OK",
				"FAIL",
			},
		},
	}

	for _, tt := range testCases {

		s := tt.setup(testPort)

		c, err := net.Dial("tcp", fmt.Sprintf(":%d", testPort))
		reader := bufio.NewReader(c)
		if err != nil {
			t.Fatalf("Test Server Failure on \n %s \n unable to talk to server on port %s \n ", tt.name, testPort)
		}
		for i, m := range tt.in {

			c.Write([]byte(m))
			buf, _, err := reader.ReadLine()
			if err != nil {
				t.Fatal("Test Server Failure on \n %s \n unable to read from server", tt.name)

			}
			raw := string(buf[:])
			if raw != tt.out[i] {
				t.Fatalf("Test Server Failure on \n %s \n Expected : %s \n Actual: %s", tt.name, tt.out[i], raw)
			}

		}

		// close client
		c.Close()

		// close server
		s.Listener.Close()

		// up test port for the next server
		testPort++

	}

}
