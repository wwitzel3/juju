// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

import (
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cert"
	"github.com/juju/juju/state"
	"github.com/juju/juju/state/watcher"
)

func init() {
	common.RegisterStandardFacade("ProcManager", 0, NewProcManagerAPI)
}

// ProcManagerAPI implements the API used by the procmanager worker.
type ProcManagerAPI struct {
	*common.EnvironWatcher

	st             *state.State
	resources      *common.Resources
	authorizer     common.Authorizer
	StateAddresser *common.StateAddresser
	canModify      bool
}

// NewProcManagerAPI creates a new instance of the ProcManager API.
func NewProcManagerAPI(st *state.State, resources *common.Resources, authorizer common.Authorizer) (*ProcManagerAPI, error) {
	if !authorizer.AuthMachineAgent() && !authorizer.AuthUnitAgent() {
		return nil, common.ErrPerm
	}
	return &ProcManagerAPI{
		EnvironWatcher: common.NewEnvironWatcher(st, resources, authorizer),
		st:             st,
		authorizer:     authorizer,
		resources:      resources,
		canModify:      authorizer.AuthEnvironManager(),
		StateAddresser: common.NewStateAddresser(st),
	}, nil
}
