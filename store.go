package main

import (
	"fmt"
	"sync"
)

type PkgDtl map[string][]string

// PkgStore is the representation of that state of the package
type PkgStore struct {

	// mutex for locking
	mutex *sync.RWMutex

	// list of that state
	Index PkgDtl
}

// Get ...
func (pkSt *PkgStore) Get(pkgName string) ([]string, bool) {

	pkSt.mutex.RLock()
	defer pkSt.mutex.RUnlock()
	pk, ok := pkSt.Index[pkgName]
	return pk, ok

}

// Remove ...
func (pkSt *PkgStore) Remove(pkgName string) bool {

	deps, ok := pkSt.Get(pkgName)

	if !ok {
		return true
	}

	fmt.Println(pkgName, deps)

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

func (pkSt *PkgStore) DepsInstalled(deps []string) bool {

	for _, v := range deps {
		_, ok := pkSt.Get(v)
		if !ok {
			return false

		}

	}
	return true

}

// Add ...
func (pkSt *PkgStore) Add(msg *Msg) bool {

	_, ok := pkSt.Get(msg.Package)

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

/*

	// check if msg has deps if len(msg.Deps) == 0 {
		pDtl := &PkgDtl{nil, nil}
		pkSt.mutex.Lock()
		pkSt.Index[msg.Package] = *pDtl
		pkSt.mutex.Unlock()
		return true

	}

	// package has dependencies

	// Does package already exist?
	pkDtl, ok := pkSt.Get(msg.Package)

	// package doesn't exist
	if !ok {
		// check deps
		chk, pkgDeps := pkSt.CheckDeps(msg.Deps)

		// dependencies exist
		// go head updates deps and insert pkg
		if chk {

			// update deps
			for k, d := range pkgDeps {

				d.ReqBy = append(d.ReqBy, msg.Package)
				pkSt.mutex.Lock()
				pkSt.Index[k] = d
				pkSt.mutex.Unlock()

			}

			// insert new package
			pDtl := &PkgDtl{msg.Deps, nil}
			pkSt.mutex.Lock()
			pkSt.Index[msg.Package] = *pDtl
			pkSt.mutex.Unlock()
			return true

		}

		return false

	}

	// pkg already exist
	// check deps
	chk, newDeps := pkSt.CheckDeps(msg.Deps)

	if chk {
		// update deps
		for k, d := range newDeps {

			d.ReqBy = append(d.ReqBy, msg.Package)
			pkSt.mutex.Lock()
			pkSt.Index[k] = d
			pkSt.mutex.Unlock()

		}

		// remove old deps
		_, oldDeps := pkSt.CheckDeps(pkDtl.Deps)

		for k, v := range oldDeps {

			for i, r := range v.ReqBy {

				// found the dep
				if r == msg.Package {
					// delete an element safely
					v.ReqBy[i] = v.ReqBy[len(v.ReqBy)-1]
					v.ReqBy[len(v.ReqBy)-1] = ""
					v.ReqBy = v.ReqBy[:len(v.ReqBy)-1]

					// lock and insert
					pkSt.mutex.Lock()
					pkSt.Index[k] = v
					pkSt.mutex.Unlock()

				}

			}

		}

		// insert new package
		pDtl := &PkgDtl{msg.Deps, nil}
		pkSt.mutex.Lock()
		pkSt.Index[msg.Package] = *pDtl
		pkSt.mutex.Unlock()
		return true

	}

	return false
*/

/*
// CheckDeps checks to see if all the dependencies are installed for a given package
func (pkSt *PkgStore) CheckDeps(deps []string) (bool, map[string]PkgDtl) {

	// return current package dependences
	deets := make(map[string]PkgDtl)

	for _, dep := range deps {
		pkgDtl, ok := pkSt.Get(dep)
		if !ok {
			return false, nil
		}

		deets[dep] = pkgDtl

	}
	return true, deets
}
*/
