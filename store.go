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
func (pkSt *PkgStore) Remove(pkgName string) bool {

	pkg, ok := pkSt.Get(pkgName)

	if !ok {
		return true
	}

	// pakcage has no dependences
	// and isn't required by anything so its okay to remove
	if len(pkg.Deps) == 0 && len(pkg.ReqBy) == 0 {

		pkSt.mutex.Lock()
		delete(pkSt.Index, pkgName)
		pkSt.mutex.Unlock()
		return true
	}

	return false

}

// Add ...
func (pkSt *PkgStore) Add(msg *Msg) bool {

	// check if msg has deps
	if len(msg.Deps) == 0 {
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

}

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
