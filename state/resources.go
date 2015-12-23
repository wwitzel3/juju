// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"io"

	"github.com/juju/errors"

	"github.com/juju/juju/resource"
)

// Resources describes the state functionality for resources.
type Resources interface {
	// ListResources returns the list of resources for the given service.
	ListResources(serviceID string) ([]resource.Resource, error)

	// SetResource stores the resource in the Juju model.
	SetResource(serviceID string, res resource.Resource, r io.Reader) error
}

var newResources func(Persistence) (Resources, error)

// SetResourcesComponent registers the function that provide the state
// functionality related to resources.
func SetResourcesComponent(fn func(Persistence) (Resources, error)) {
	newResources = fn
}

// Resources returns the resources functionality for the current state.
func (st *State) Resources() (Resources, error) {
	if newResources == nil {
		return nil, errors.Errorf("resources not supported")
	}

	persist := st.newPersistence()
	resources, err := newResources(persist)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return resources, nil
}
