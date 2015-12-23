// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api

// TODO(ericsnow) Eliminate the dependence on apiserver if possible.

import (
	"github.com/juju/errors"
	"github.com/juju/names"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"

	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/resource"
)

// ServiceTag2ID converts the provided tag into a service ID.
func ServiceTag2ID(tagStr string) (string, error) {
	kind, err := names.TagKind(tagStr)
	if err != nil {
		return "", errors.Annotatef(err, "could not determine tag type from %q", tagStr)
	}
	if kind != names.ServiceTagKind {
		return "", errors.Errorf("expected service tag, got %q", tagStr)
	}

	tag, err := names.ParseTag(tagStr)
	if err != nil {
		return "", errors.Errorf("invalid service tag %q", tagStr)
	}
	return tag.Id(), nil
}

// Resource2API converts a resource.Resource into
// a Resource struct.
func Resource2API(res resource.Resource) Resource {
	return Resource{
		CharmResource: CharmResource2API(res.Resource),
		Username:      res.Username,
		Timestamp:     res.Timestamp,
	}
}

// APIResult2Resources converts a ResourcesResult into []resource.Resource.
func APIResult2Resources(apiResult ResourcesResult) ([]resource.Resource, error) {
	var result []resource.Resource

	if apiResult.Error != nil {
		// TODO(ericsnow) Return the resources too?
		err, _ := common.RestoreError(apiResult.Error)
		return nil, errors.Trace(err)
	}

	for _, apiRes := range apiResult.Resources {
		res, err := API2Resource(apiRes)
		if err != nil {
			// This could happen if the server is misbehaving
			// or non-conforming.
			// TODO(ericsnow) Aggregate errors?
			return nil, errors.Annotate(err, "got bad data from server")
		}
		result = append(result, res)
	}

	return result, nil
}

// API2Resource converts an API Resource struct into
// a resource.Resource.
func API2Resource(apiRes Resource) (resource.Resource, error) {
	var res resource.Resource

	charmRes, err := API2CharmResource(apiRes.CharmResource)
	if err != nil {
		return res, errors.Trace(err)
	}

	res = resource.Resource{
		Resource:  charmRes,
		Username:  apiRes.Username,
		Timestamp: apiRes.Timestamp,
	}

	if err := res.Validate(); err != nil {
		return res, errors.Trace(err)
	}

	return res, nil
}

// CharmResource2API converts a charm resource into
// a CharmResource struct.
func CharmResource2API(res charmresource.Resource) CharmResource {
	return CharmResource{
		Name:        res.Name,
		Type:        res.Type.String(),
		Path:        res.Path,
		Comment:     res.Comment,
		Origin:      res.Origin.String(),
		Revision:    res.Revision,
		Fingerprint: res.Fingerprint.Bytes(),
	}
}

// API2CharmResource converts an API CharmResource struct into
// a charm resource.
func API2CharmResource(apiInfo CharmResource) (charmresource.Resource, error) {
	var res charmresource.Resource

	rtype, err := charmresource.ParseType(apiInfo.Type)
	if err != nil {
		return res, errors.Trace(err)
	}

	origin, err := charmresource.ParseOrigin(apiInfo.Origin)
	if err != nil {
		return res, errors.Trace(err)
	}

	fp, err := charmresource.NewFingerprint(apiInfo.Fingerprint)
	if err != nil {
		return res, errors.Trace(err)
	}
	if err := fp.Validate(); err != nil {
		return res, errors.Trace(err)
	}

	res = charmresource.Resource{
		Meta: charmresource.Meta{
			Name:    apiInfo.Name,
			Type:    rtype,
			Path:    apiInfo.Path,
			Comment: apiInfo.Comment,
		},
		Origin:      origin,
		Revision:    apiInfo.Revision,
		Fingerprint: fp,
	}

	if err := res.Validate(); err != nil {
		return res, errors.Trace(err)
	}
	return res, nil
}
