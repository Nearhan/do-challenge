package main

import (
	"reflect"
	"testing"
)

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
