// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"strconv"

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
	info, ok := data.(*procmanager.ProcessInfo)
	if !ok {
		return errors.Errorf("invalid data type %T", data)
	}
	pDoc, ok := doc.(*processInfoDoc)
	if !ok {
		return errors.Errorf("invalid data type %T", doc)
	}

	status, err := strconv.Atoi(info.Details.Status)
	if err != nil {
		status = -1 // unknown
	}

	*pDoc = processInfoDoc{
		Image:      info.Image,
		Args:       info.Args,
		Desc:       info.Desc,
		Plugin:     info.Plugin,
		Storage:    info.Storage,
		Networking: info.Networking,
		UniqueID:   info.Details.UniqueID,
		Status:     state.Life(status),
	}
	return nil
}

// ConvertOut implements state.CollectionHandler
func (handler) ConvertOut(data, doc interface{}) error {
	info, ok := data.(*procmanager.ProcessInfo)
	if !ok {
		return errors.Errorf("invalid data type %T", data)
	}
	pDoc, ok := doc.(*processInfoDoc)
	if !ok {
		return errors.Errorf("invalid data type %T", doc)
	}

	*info = procmanager.ProcessInfo{
		Image:      pDoc.Image,
		Args:       pDoc.Args,
		Desc:       pDoc.Desc,
		Plugin:     pDoc.Plugin,
		Storage:    pDoc.Storage,
		Networking: pDoc.Networking,
		Details: procmanager.ProcessDetails{
			UniqueID: pDoc.UniqueID,
			Status:   pDoc.Status.String(),
		},
	}
	return nil
}

// NewId implements state.CollectionHandler
func (handler) NewId(envUUID string, data interface{}) (string, error) {
	info, ok := data.(*procmanager.ProcessInfo)
	if !ok {
		return "", errors.Errorf("invalid data type %T", data)
	}
	return info.Details.UniqueID + "." + envUUID, nil
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
