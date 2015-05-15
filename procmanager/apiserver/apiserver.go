// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common"
	//"github.com/juju/juju/procmanager"
	"github.com/juju/juju/procmanager/apiclient"
	"github.com/juju/juju/state"
)

func init() {
	common.RegisterStandardFacade("ProcManager", 1, NewAPI)
}

// API serves procmanager-specific API methods.
type API struct {
	st *state.State
}

// NewAPI creates a new instance of the ProcManager API facade.
func NewAPI(st *state.State, resources *common.Resources, authorizer common.Authorizer) (*API, error) {
	if !authorizer.AuthClient() {
		return nil, errors.Trace(common.ErrPerm)
	}
	// Build the API.
	b := API{
		st: st,
	}
	return &b, nil
}

func (a *API) Add(args apiclient.AddArgs) string {
	return ""
}

func (a *API) Remove() string {
	return ""
}

func (a *API) Info() string {
	return ""
}
