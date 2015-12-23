// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// The resource package provides the functionality of the "resources"
// feature in Juju.
package resource

import (
	"time"

	"github.com/juju/errors"
	"gopkg.in/juju/charm.v6-unstable/resource"
)

// Resource defines a single resource within Juju state.
type Resource struct {
	resource.Resource

	// Username is the ID of the user that added the revision
	// to the model (whether implicitly or explicitly).
	Username string

	// Timestamp indicates when the resource was added to the model.
	Timestamp time.Time
}

// Validate ensures that the spec is valid.
func (res Resource) Validate() error {
	if err := res.Resource.Validate(); err != nil {
		return errors.Annotate(err, "bad info")
	}

	// TODO(ericsnow) Require that Username be set if timestamp is?

	if res.Timestamp.IsZero() {
		if res.Username != "" {
			return errors.NewNotValid(nil, "missing timestamp")
		}
	}

	return nil
}
