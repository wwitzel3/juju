// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"github.com/juju/errors"

	"github.com/juju/juju/procmanager"
	"github.com/juju/juju/state"
)

// TODO(ericsnow) Factor in envUUID?

// processInfo holds information about a process.
type processInfoDoc struct {
	DocID   string `bson:"_id"`
	EnvUUID string `bson:"env-uuid"`

	Image      string     `bson:"image"`
	Args       string     `bson:"args"`
	Desc       string     `bson:"desc"`
	Plugin     string     `bson:"plugin"`
	Storage    string     `bson:"storage"`
	Networking string     `bson:"networking"`
	UniqueID   string     `bson:"unique-id"`
	Status     state.Life `bson:"lift"`
}

func init() {
	state.RegisterCollectionHandler(state.ProcessesC, &handler{})
}

type handler struct{}

// ConvertIn implements state.CollectionHandler
func (handler) ConvertIn(doc, data interface{}) error {
	return nil
}

// ConvertOut implements state.CollectionHandler
func (handler) ConvertOut(data, doc interface{}) error {
	return nil
}

// NewId implements state.CollectionHandler
func (handler) NewId(envUUID string, data interface{}) (string, error) {
	return "", nil
}

// Register adds info about a process to Juju's state.
func Register(st *state.State, info *procmanager.ProcessInfo) (string, error) {
	return st.InsertOne(state.ProcessesC, st.EnvironUUID(), info)
}

// Unregister removes info about a process from Juju's state.
func Unregister(st *state.State, uuid string) error {
	return st.RemoveOne(state.ProcessesC, uuid)
}

// Info retrieves info about a process from Juju's state.
func Info(st *state.State, uuid string) (*procmanager.ProcessInfo, error) {
	var result procmanager.ProcessInfo

	if err := st.GetOne(state.ProcessesC, uuid, &result); err != nil {
		return nil, errors.Trace(err)
	}

	return &result, nil
}
