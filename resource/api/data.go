// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api

import (
	"io"
	"time"

	"github.com/juju/errors"
	"github.com/juju/names"

	"github.com/juju/juju/apiserver/params"
)

// Resource contains info about a Resource.
type Resource struct {
	CharmResource

	// Username is the ID of the user that added the revision
	// to the model (whether implicitly or explicitly).
	Username string `json:"username"`

	// Timestamp indicates when the resource was added to the model.
	Timestamp time.Time `json:"timestamp"`
}

// CharmResource contains the definition for a resource.
type CharmResource struct {
	// Name identifies the resource.
	Name string `json:"name"`

	// Type is the name of the resource type.
	Type string `json:"type"`

	// Path is where the resource will be stored.
	Path string `json:"path"`

	// Comment contains user-facing info about the resource.
	Comment string `json:"comment,omitempty"`

	// Origin is where the resource will come from.
	Origin string `json:"origin"`

	// Revision is the revision, if applicable.
	Revision int `json:"revision"`

	// Fingerprint is the SHA-384 checksum for the resource blob.
	Fingerprint []byte `json:"fingerprint"`
}

type UploadEntity struct {
	Tag  string
	Name string
	Blob io.Reader
}

type UploadArgs struct {
	Entities []UploadEntity
}

// NewUploadArgs returns the arguments for the Upload endpoint.
func NewUploadArgs(service string, name string, blob io.Reader) (UploadArgs, error) {
	var args UploadArgs
	if !names.IsValidService(service) {
		return args, errors.Errorf("invalid service %q", service)
	}

	args.Entities = append(args.Entities, UploadEntity{
		Tag:  names.NewServiceTag(service).String(),
		Name: name,
		Blob: blob,
	})
	return args, nil
}

type UploadResults struct {
	Results []UploadResult
}

type UploadResult struct {
	params.ErrorResult
}

func NewUploadResult(tag string, name string, blob io.Reader) (UploadResult, error) {
	return UploadResult{}, nil
}
