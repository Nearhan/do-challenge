package main

import (
	"sync"
)

// PkgDtl ...
type PkgDtl map[string][]string

// PkgStore is the representation of that state of the package
type PkgStore struct {

	// mutex for locking
	mutex *sync.RWMutex

	// list of that state
	Index PkgDtl
}

// Get ...
func (pkSt *PkgStore) Get(pkgName string) bool {

	pkSt.mutex.RLock()
	defer pkSt.mutex.RUnlock()
	_, ok := pkSt.Index[pkgName]
	return ok

}

// Remove ...
func (pkSt *PkgStore) Remove(pkgName string) bool {

	ok := pkSt.Get(pkgName)

	if !ok {
		return true
	}

	//iterate over everything
	if !pkSt.hasDependencies(pkgName) {

		pkSt.mutex.Lock()
		delete(pkSt.Index, pkgName)
		pkSt.mutex.Unlock()
		return true
	}

	return false

}

func (pkSt *PkgStore) hasDependencies(pkgName string) bool {

	pkSt.mutex.RLock()
	defer pkSt.mutex.RUnlock()

	for _, v := range pkSt.Index {
		for _, d := range v {
			if d == pkgName {
				return true
			}
		}

	}

	return false

}

// DepsInstalled ...
func (pkSt *PkgStore) DepsInstalled(deps []string) bool {

	for _, v := range deps {
		ok := pkSt.Get(v)
		if !ok {
			return false

		}

	}
	return true

}

// Add ...
func (pkSt *PkgStore) Add(msg *Msg) bool {

	ok := pkSt.Get(msg.Package)

	if ok {

		if len(msg.Deps) > 0 {
			if pkSt.DepsInstalled(msg.Deps) {

				pkSt.mutex.Lock()
				pkSt.Index[msg.Package] = msg.Deps
				pkSt.mutex.Unlock()
				return true
			}
			return false
		}

		pkSt.mutex.Lock()
		pkSt.Index[msg.Package] = msg.Deps
		pkSt.mutex.Unlock()
		return true

	}
	if len(msg.Deps) > 0 {
		if pkSt.DepsInstalled(msg.Deps) {

			pkSt.mutex.Lock()
			pkSt.Index[msg.Package] = msg.Deps
			pkSt.mutex.Unlock()
			return true
		}
		return false
	}

	pkSt.mutex.Lock()
	pkSt.Index[msg.Package] = msg.Deps
	pkSt.mutex.Unlock()
	return true

}
