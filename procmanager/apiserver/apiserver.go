// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver

import (
	"github.com/juju/errors"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/procmanager"
	"github.com/juju/juju/procmanager/apiclient"
	procstate "github.com/juju/juju/procmanager/state"
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

func (a *API) Add(args apiclient.ProcessInfo) (*apiclient.ProcessID, error) {
	info := procmanager.ProcessInfo{
		Image:      args.Image,
		Args:       args.Args,
		Desc:       args.Desc,
		Plugin:     args.Plugin,
		Storage:    args.Storage,
		Networking: args.Networking,
		Details: procmanager.ProcessDetails{
			UniqueID: args.UniqueID,
			Status:   args.Status,
		},
	}
	uuid, err := procstate.Register(a.st, info)
	if err != nil {
		return nil, errors.Trace(err)
	}
	result := apiclient.ProcessID{
		UUID: uuid,
	}
	return &result, nil
}

func (a *API) Remove(args apiclient.ProcessID) error {
	err := procstate.Unregister(a.st, args.UUID)
	return errors.Trace(err)
}

func (a *API) Info(args apiclient.ProcessID) (*apiclient.ProcessInfo, error) {
	info, err := procstate.Info(a.st, args.UUID)
	if err != nil {
		return nil, errors.Trace(err)
	}

	result := apiclient.ProcessInfo{
		Image:      info.Image,
		Args:       info.Args,
		Desc:       info.Desc,
		Plugin:     info.Plugin,
		Storage:    info.Storage,
		Networking: info.Networking,
		UniqueID:   info.Details.UniqueID,
		Status:     info.Details.Status,
	}
	return &result, nil
}
