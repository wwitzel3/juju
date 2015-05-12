// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

import (
	"github.com/juju/juju/api/base"
)

const procManagerAPI = "ProcManager"

// State provides access to the Rsyslog API facade.
type State struct {
	facade base.FacadeCaller
}

// NewState creates a new client-side Rsyslog facade.
func NewState(caller base.APICaller) *State {
	return &State{facade: base.NewFacadeCaller(caller, procManagerAPI)}
}

func (st *State) Launch() error {
	return nil
}

func (st *State) Destroy() error {
	return nil
}
