package main

import "sync"

// PkgDetail ...
type PkgDtl struct {

	// Package dependencies
	Deps []string

	// What packages are required by this package
	ReqBy []string
}

// PkgStore is the representation of that state of the package
type PkgStore struct {

	// mutex for locking
	mutex *sync.RWMutex

	// list of that state
	Index map[string]PkgDtl
}

// Get ...
func (pkSt *PkgStore) Get(pkgName string) (PkgDtl, bool) {

	pkSt.mutex.RLock()
	defer pkSt.mutex.RUnlock()
	pk, ok := pkSt.Index[pkgName]
	return pk, ok

}

// Remove ...
func (pkSt *PkgStore) Remove(pkgName string) {

	pkSt.mutex.Lock()
	delete(pkSt.Index, pkgName)
	pkSt.mutex.Unlock()

}

// Add ...
func (pkSt *PkgStore) Add(msg *Msg) {

	pDtl := &PkgDtl{msg.Deps, nil}
	pkSt.mutex.Lock()
	pkSt.Index[msg.Package] = *pDtl
	pkSt.mutex.Unlock()

}

// CheckDeps checks to see if all the dependencies are installed for a given package
func (pkSt *PkgStore) CheckDeps(deps []string) bool {

	for _, dep := range deps {
		_, ok := pkSt.Get(dep)
		if !ok {
			return false
		}
	}
	return true
}
