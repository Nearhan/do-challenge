package main

import (
	"reflect"
	"sync"
	"testing"
)

func compareStore(a, b *PkgStore) bool {

	if len(a.Index) != len(b.Index) {
		return false
	}

	for k, v := range a.Index {
		d, ok := b.Index[k]
		if !ok {
			return false
		}
		if !compareSlices(v.Deps, d.Deps) || !compareSlices(v.ReqBy, d.ReqBy) {
			return false
		}
	}

	return true
}

func compareSlices(c, d []string) bool {
	return reflect.DeepEqual(c, d)
}

func makeStore() *PkgStore {

	return &PkgStore{&sync.RWMutex{}, make(map[string]PkgDtl)}
}

func makeStoreWithState(m map[string]PkgDtl) *PkgStore {
	s := makeStore()
	s.Index = m
	return s

}

func complexState1() *PkgStore {

	//&Msg{"INDEX", "git", []string{"ubuntu", firefox", "opera"}},
	return makeStoreWithState(map[string]PkgDtl{
		"git":     PkgDtl{[]string{"osx", "ubuntu", "windows"}, nil},
		"osx":     PkgDtl{make([]string, 0), []string{"git"}},
		"ubuntu":  PkgDtl{make([]string, 0), []string{"git"}},
		"windows": PkgDtl{make([]string, 0), []string{"git"}},
		"firefox": PkgDtl{make([]string, 0), []string{"cat"}},
		"opera":   PkgDtl{make([]string, 0), []string{"dog"}},
	})
}

func complexState2() *PkgStore {

	//&Msg{"INDEX", "git", []string{"ubuntu", firefox", "opera"}},
	return makeStoreWithState(map[string]PkgDtl{
		"git":     PkgDtl{[]string{"ubuntu", "firefox", "opera"}, make([]string, 0)},
		"osx":     PkgDtl{make([]string, 0), make([]string, 0)},
		"ubuntu":  PkgDtl{make([]string, 0), []string{"git"}},
		"windows": PkgDtl{make([]string, 0), make([]string, 0)},
		"firefox": PkgDtl{make([]string, 0), []string{"git", "cat"}},
		"opera":   PkgDtl{make([]string, 0), []string{"git", "dog"}},
	})
}

func TestStoreRemove(t *testing.T) {
	testStates := []struct {
		in        string    // input
		initStore *PkgStore // inital store state
		finStore  *PkgStore // final store state
		ret       bool      // Remove return value
	}{
		{
			"git",
			makeStoreWithState(map[string]PkgDtl{"git": PkgDtl{nil, nil}}),
			makeStore(),
			true,
		},
		{
			"git",
			makeStoreWithState(map[string]PkgDtl{"git": PkgDtl{[]string{"osx"}, nil}}),
			makeStoreWithState(map[string]PkgDtl{"git": PkgDtl{[]string{"osx"}, nil}}),
			false,
		},
		{
			"vim",
			makeStoreWithState(map[string]PkgDtl{"vim": PkgDtl{[]string{"git", "go", "python"}, nil}}),
			makeStoreWithState(map[string]PkgDtl{"vim": PkgDtl{[]string{"git", "go", "python"}, nil}}),
			false,
		},
	}

	for _, tt := range testStates {
		out := tt.initStore.Remove(tt.in)

		if !reflect.DeepEqual(tt.initStore, tt.finStore) || out != tt.ret {
			t.Fatalf(" Test PkgStore Remove failure | Expected: %s | Acutal: %s ", tt.finStore, tt.initStore)
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
			"Index with empty store and no dependencies",
			&Msg{"INDEX", "git", nil},
			makeStore(),
			makeStoreWithState(map[string]PkgDtl{"git": PkgDtl{nil, nil}}),
			true,
		},
		{
			"Index with empty store with dependencies",
			&Msg{"INDEX", "git", []string{"vim", "osx"}},
			makeStore(),
			makeStore(),
			false,
		},
		{

			"Index with stateful store and dependencies",
			&Msg{"INDEX", "git", []string{"vim", "osx"}},
			makeStoreWithState(map[string]PkgDtl{"vim": PkgDtl{nil, nil}, "osx": PkgDtl{nil, nil}}),
			makeStoreWithState(map[string]PkgDtl{
				"vim": PkgDtl{nil, []string{"git"}},
				"osx": PkgDtl{nil, []string{"git"}},
				"git": PkgDtl{[]string{"vim", "osx"}, nil}},
			),
			true,
		},
		{
			"Index with stateful store swap dependencies",
			&Msg{"INDEX", "git", []string{"vim"}},
			makeStoreWithState(map[string]PkgDtl{
				"vim": PkgDtl{make([]string, 0), make([]string, 0)},
				"osx": PkgDtl{make([]string, 0), []string{"git"}},
				"git": PkgDtl{[]string{"osx"}, make([]string, 0)}},
			),
			makeStoreWithState(map[string]PkgDtl{
				"vim": PkgDtl{make([]string, 0), []string{"git"}},
				"osx": PkgDtl{make([]string, 0), make([]string, 0)},
				"git": PkgDtl{[]string{"vim"}, make([]string, 0)}},
			),
			true,
		},
		{
			"Index with complex types",
			&Msg{"INDEX", "git", []string{"ubuntu", "firefox", "opera"}},
			complexState1(),
			complexState2(),
			true,
		},
	}

	for _, tt := range testStates {

		out := tt.initStore.Add(tt.in)
		//fmt.Println(tt.initStore)

		//fmt.Printf("Test PkgStore Add failure on: %s \n Expected: %s \n Acutal: %s \n RETURN %s %s", tt.name, tt.finStore, tt.initStore, tt.ret, out)

		if !reflect.DeepEqual(tt.initStore.Index["git"], tt.finStore.Index["git"]) || out != tt.ret {
			t.Fatalf("Test PkgStore Add failure on: %s \n Expected: %s \n Acutal: %s \n RETURN %s %s", tt.name, tt.finStore, tt.initStore, tt.ret, out)
		}

	}

}

func TestStoreGet(t *testing.T) {

	testStates := []struct {
		in        string    // input
		initStore *PkgStore // inital store state
		ret1      PkgDtl    // final store state
		ret2      bool      // Index return value
	}{
		{
			"git",
			makeStore(),
			PkgDtl{nil, nil},
			false,
		},
		{
			"git",
			makeStoreWithState(map[string]PkgDtl{"git": PkgDtl{nil, nil}}),
			PkgDtl{nil, nil},
			true,
		},
	}

	for _, tt := range testStates {
		out1, out2 := tt.initStore.Get(tt.in)

		if !reflect.DeepEqual(out1, tt.ret1) || out2 != tt.ret2 {
			t.Fatalf("Test PkgStore Query failure Expected: %s  & %s \n Acutal: %s & %s \n", tt.ret1, tt.ret2, out1, out2)
		}

	}

}
