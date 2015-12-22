// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"io"
	"path"

	"github.com/juju/errors"

	"github.com/juju/juju/resource"
)

type resourcePersistence interface {
	// ListResources returns the resource data for the given service ID.
	ListResources(serviceID string) ([]resource.Resource, error)

	// SetStagedResource adds the resource in a separate staging area
	// if the resource isn't already staged. If the resource already
	// exists then it is treated as unavailable as long as the new one
	// is staged.
	SetStagedResource(id, serviceID string, res resource.Resource) error

	// UnstageResource ensures that the resource is removed
	// from the staging area. If it isn't in the staging area
	// then this is a noop.
	UnstageResource(id, serviceID string) error

	// SetResource stores the resource info. If the resource
	// is already staged then it is unstaged.
	SetResource(id, serviceID string, res resource.Resource) error
}

type resourceStorage interface {
	// Put stores the content of the reader into the storage.
	Put(path, hash string, length int64, r io.Reader) error

	// Delete removes the identified data from the storage.
	Delete(path string) error
}

type resourceState struct {
	persist resourcePersistence
	storage resourceStorage
}

// ListResources returns the resource data for the given service ID.
func (st resourceState) ListResources(serviceID string) ([]resource.Resource, error) {
	resources, err := st.persist.ListResources(serviceID)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return resources, nil
}

// TODO(ericsnow) Separate setting the metadata from storing the blob?

// SetResource stores the resource in the Juju model.
func (st resourceState) SetResource(serviceID string, res resource.Resource, r io.Reader) error {
	if err := res.Validate(); err != nil {
		return errors.Annotate(err, "bad resource metadata")
	}
	id := res.Name
	hash := string(res.Fingerprint.Bytes())

	// TODO(ericsnow) Do something else if r is nil?

	// We use a staging approach for adding the resource metadata
	// to the model. This is necessary because the resource data
	// is stored separately and adding to both should be an atomic
	// operation.

	if err := st.persist.SetStagedResource(id, serviceID, res); err != nil {
		return errors.Trace(err)
	}

	path := storagePath(res.Name, serviceID)
	if err := st.storage.Put(path, hash, res.Size, r); err != nil {
		if err := st.persist.UnstageResource(id, serviceID); err != nil {
			logger.Errorf("could not unstage resource %q (service %q): %v", res.Name, serviceID, err)
		}
		return errors.Trace(err)
	}

	if err := st.persist.SetResource(id, serviceID, res); err != nil {
		if err := st.storage.Delete(path); err != nil {
			logger.Errorf("could not remove resource %q (service %q) from storage: %v", res.Name, serviceID, err)
		}
		if err := st.persist.UnstageResource(id, serviceID); err != nil {
			logger.Errorf("could not unstage resource %q (service %q): %v", res.Name, serviceID, err)
		}
		return errors.Trace(err)
	}

	return nil
}

func storagePath(id, serviceID string) string {
	return path.Join("service-"+serviceID, "resources", id)
}
