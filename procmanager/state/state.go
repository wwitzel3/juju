// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"github.com/juju/juju/procmanager"
	"github.com/juju/juju/state"
)

// TODO(ericsnow) Factor in envUUID?

// Register adds info about a process to Juju's state.
func Register(st *state.State, info procmanager.ProcessInfo) (string, error) {
	// TODO(ericsnow) finish
	return "", nil
}

// Unregister removes info about a process from Juju's state.
func Unregister(st *state.State, uuid string) error {
	// TODO(ericsnow) finish
	return nil
}

// Info retrieves info about a process from Juju's state.
func Info(st *state.State, uuid string) (*procmanager.ProcessInfo, error) {
	var result procmanager.ProcessInfo

	// TODO(ericsnow) finish
	return &result, nil
}
