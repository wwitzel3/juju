// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

// UnitProcess represents the state of a long running process
// on a unit.
type UnitProcess interface {
	// Owner returns the tag of the service or unit that owns this process.
	Owner() names.Tag

	// Life reports whether the process is Alive, Dying, or Dead.
	Life() Life

	// CharmURL returns the charm URL taht created this process.
	CharmURL() *charm.URL
}

type unitProcess struct {
	st  *State
	doc unitProcessDoc
}

func (u unitProcess) Owner() names.Tag {
	tag, err := names.ParseTag(s.doc.Owner)
	if err != nil {
		panic(err)
	}
	return tag
}

func (u unitProcess) Life() Life {
	return s.doc.Life
}

func (u unitProcess) CharmURL() *charm.URL {
	return s.doc.CharmURL
}
