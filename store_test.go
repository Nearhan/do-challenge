package main

import (
	"reflect"
	"sync"
	"testing"
)

func makeStore() *PkgStore {

	return &PkgStore{&sync.RWMutex{}, make(PkgDtl)}
}

func makeStoreWithState(m *PkgDtl) *PkgStore {
	s := makeStore()
	s.Index = *m
	return s

}

func makeEmptySliceStr() []string {

	return make([]string, 0)

}

func TestStoreRemove(t *testing.T) {
	testStates := []struct {
		name      string    // name of the test
		in        string    // input
		initStore *PkgStore // inital store state
		finStore  *PkgStore // final store state
		ret       bool      // Remove return value
	}{
		{
			"Should Remove with empty store and no deps",
			"git",
			makeStore(),
			makeStore(),
			true,
		},
		{
			"Should Not Remove with stateful store and deps",
			"git",
			makeStoreWithState(&PkgDtl{"git": []string{}, "osx": []string{"git"}}),
			makeStoreWithState(&PkgDtl{"git": []string{}, "osx": []string{"git"}}),
			false,
		},
		{

			"Should Remove with stateful store and deps",
			"vim",
			makeStoreWithState(&PkgDtl{"vim": []string{}}),
			makeStore(),
			true,
		},
	}

	for _, tt := range testStates {
		out := tt.initStore.Remove(tt.in)

		if !reflect.DeepEqual(tt.initStore, tt.finStore) || out != tt.ret {
			t.Fatalf(" Test PkgStore Remove failure on %s \n | Expected: %s \n | Acutal: %s \n", tt.name, tt.finStore, tt.initStore)
		}

	}

}

func TestStoreAdd(t *testing.T) {
	testStates := []struct {
		name      string    // name of the type of check we are doing
		in        *Msg      // input
		initStore *PkgStore // inital store state
		finStore  *PkgStore // final store state
		ret       bool      // Index return value
	}{
		{
			"Should Index with empty store and no dependencies",
			&Msg{"INDEX", "git", make([]string, 0)},
			makeStore(),
			makeStoreWithState((&PkgDtl{"git": make([]string, 0)})),
			true,
		},
		{
			"Should NOT Index with empty store with dependencies",
			&Msg{"INDEX", "git", []string{"vim", "osx"}},
			makeStore(),
			makeStore(),
			false,
		},
		{

			"Should Index with stateful store and dependencies",
			&Msg{"INDEX", "git", []string{"vim", "osx"}},
			makeStoreWithState(&PkgDtl{
				"vim": makeEmptySliceStr(),
				"osx": makeEmptySliceStr(),
			}),
			makeStoreWithState(&PkgDtl{
				"vim": makeEmptySliceStr(),
				"osx": makeEmptySliceStr(),
				"git": []string{"vim", "osx"},
			}),
			true,
		},
		{
			"Should Index with stateful store swap dependencies",
			&Msg{"INDEX", "git", []string{"vim"}},
			makeStoreWithState(&PkgDtl{
				"vim": makeEmptySliceStr(),
				"osx": makeEmptySliceStr(),
				"git": []string{"osx"},
			}),
			makeStoreWithState(&PkgDtl{
				"vim": makeEmptySliceStr(),
				"osx": makeEmptySliceStr(),
				"git": []string{"vim"},
			}),
			true,
		},
	}

	for _, tt := range testStates {

		out := tt.initStore.Add(tt.in)

		if !reflect.DeepEqual(tt.initStore.Index["git"], tt.finStore.Index["git"]) || out != tt.ret {
			t.Fatalf("Test PkgStore Add failure on: %s \n Expected: %s \n Acutal: %s \n RETURN %s %s", tt.name, tt.finStore, tt.initStore, tt.ret, out)
		}

	}

}

func TestStoreGet(t *testing.T) {

	testStates := []struct {
		name      string    // name of test case
		in        string    // input
		initStore *PkgStore // inital store state
		ret       bool      // Index return value
	}{
		{
			"GET return false with package missing",
			"git",
			makeStore(),
			false,
		},
		{

			"GET return true with package missing",
			"git",
			makeStoreWithState(&PkgDtl{"git": makeEmptySliceStr()}),
			true,
		},
	}

	for _, tt := range testStates {
		out := tt.initStore.Get(tt.in)

		if out != tt.ret {
			t.Fatalf("Test PkgStore GET failure on %s \n Expected: %s \n Acutal: %s \n", tt.name, tt.ret, out)
		}

	}

}
