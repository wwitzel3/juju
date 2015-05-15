// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiclient

import (
	"github.com/juju/errors"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/procmanager"
)

type apiState interface {
	base.APICallCloser
}

// Client wraps the backups API for the client.
type Client struct {
	base.ClientFacade
	facade     base.FacadeCaller
	baseFacade base.FacadeCaller
}

// NewClient returns a new procmanager client.
func NewClient(st apiState) *Client {
	frontend, backend := base.NewClientFacade(st, "ProcManager")
	return &Client{
		ClientFacade: frontend,
		facade:       backend,
	}
}

// Add adds information about a launch process to state.
func (c *Client) Add(info procmanager.ProcessInfo) (string, error) {
	var result ProcessID
	args := ProcessInfo{
		Image:      info.Image,
		Args:       info.Args,
		Desc:       info.Desc,
		Plugin:     info.Plugin,
		Storage:    info.Storage,
		Networking: info.Networking,
		UniqueID:   info.Details.UniqueID,
		Status:     info.Details.Status,
	}
	if err := c.facade.FacadeCall("Add", args, &result); err != nil {
		return "", errors.Trace(err)
	}
	return result.UUID, nil
}

// Remove removes information about a launch process to state.
func (c *Client) Remove(uuid string) error {
	args := ProcessID{UUID: uuid}
	if err := c.facade.FacadeCall("Remove", args, nil); err != nil {
		return errors.Trace(err)
	}
	return nil
}

// Info retrieves the information about a launch process from state.
func (c *Client) Info(uuid string) (*procmanager.ProcessInfo, error) {
	var result ProcessInfo
	args := ProcessID{UUID: uuid}
	if err := c.facade.FacadeCall("Info", args, &result); err != nil {
		return nil, errors.Trace(err)
	}
	info := procmanager.ProcessInfo{
		Image:      result.Image,
		Args:       result.Args,
		Desc:       result.Desc,
		Plugin:     result.Plugin,
		Storage:    result.Storage,
		Networking: result.Networking,
		Details: procmanager.ProcessDetails{
			UniqueID: result.UniqueID,
			Status:   result.Status,
		},
	}
	return &info, nil
}
