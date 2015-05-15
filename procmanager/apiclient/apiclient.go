// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiclient

import (
	"github.com/juju/errors"
	"github.com/juju/loggo"

	"github.com/juju/juju/api/base"
	"github.com/juju/juju/procmanager"
)

var logger = loggo.GetLogger("juju.procmanager.api")

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
	var result string
	args := AddArgs{
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
	return result, nil
}

// Remove removes information about a launch process to state.
func (c *Client) Remove() string {
	// make facadecall call with params, return params
	return ""
}

// Info retrieves the information about a launch process from state.
func (c *Client) Info() string {
	// make facadecall call with params, return params
	return ""
}
