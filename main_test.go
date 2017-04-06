package main

import (
	"os"
	"testing"
)

// TestMain Entry into testing suite
func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup() {

}

func tearDown() {

}
