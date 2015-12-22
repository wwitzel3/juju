// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"
)

var logger = loggo.GetLogger("juju.resource.state")

// Persistence is the state persistence functionality needed for resources.
type Persistence interface {
	resourcePersistence
}

// Storage is the state storage functionality needed for resources.
type Storage interface {
	resourceStorage
}

// RawState defines the functionality needed from state.State for resources.
type RawState interface {
	// Persistence exposes the state data persistence needed for resources.
	Persistence() (Persistence, error)

	// Storage exposes the state blob storage needed for resources.
	Storage() (Storage, error)
}

// State exposes the state functionality needed for resources.
type State struct {
	*resourceState
}

// NewState returns a new State for the given raw Juju state.
func NewState(raw RawState) (*State, error) {
	logger.Tracef("wrapping state for resources")

	persist, err := raw.Persistence()
	if err != nil {
		return nil, errors.Trace(err)
	}

	storage, err := raw.Storage()
	if err != nil {
		return nil, errors.Trace(err)
	}

	st := &State{
		resourceState: &resourceState{
			persist,
			storage,
		},
	}
	return st, nil
}
