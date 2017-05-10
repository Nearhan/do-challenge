package main

import (
	"sync"
)

// PkgDtl type alias
type PkgDtl map[string][]string

// PkgStore is the representation of that state of the package
type PkgStore struct {

	// mutex for locking
	mutex *sync.RWMutex

	// list of that state
	Index PkgDtl
}

// Get tries to get a package
func (pkSt *PkgStore) Get(pkgName string) bool {

	pkSt.mutex.RLock()
	defer pkSt.mutex.RUnlock()

	_, ok := pkSt.Index[pkgName]
	return ok

}

// Remove tries to remove a package if no other
// packages depend on it
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

// checks to see any package depends on pkgName
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

// DepsInstalled checks to see if the packages inside deps are all installed
func (pkSt *PkgStore) DepsInstalled(deps []string) bool {

	for _, v := range deps {
		ok := pkSt.Get(v)
		if !ok {
			return false

		}

	}
	return true

}

// Add tries to add a package to the store
// checks its deps first to see if they are installed
func (pkSt *PkgStore) Add(msg *Msg) bool {

	// does it have deps?
	if len(msg.Deps) > 0 {

		// check if its deps are installed
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
