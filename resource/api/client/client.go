// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client

import (
	"io"
	"net/http"

	"github.com/juju/errors"
	"github.com/juju/loggo"

	"github.com/juju/juju/apiserver/httpattachment"
	"github.com/juju/juju/resource"
	"github.com/juju/juju/resource/api"
)

var logger = loggo.GetLogger("juju.resource.api.client")

// TODO(ericsnow) Move FacadeCaller to a component-central package.

// FacadeCaller has the api/base.FacadeCaller methods needed for the component.
type FacadeCaller interface {
	FacadeCall(request string, params, response interface{}) error
}

// Doer
type Doer interface {
	Do(req *http.Request, body io.ReadSeeker, resp interface{}) error
}

// Client is the public client for the resources API facade.
type Client struct {
	FacadeCaller
	io.Closer
	doer Doer
}

// NewClient returns a new Client for the given raw API caller.
func NewClient(caller FacadeCaller, doer Doer, closer io.Closer) *Client {
	return &Client{
		FacadeCaller: caller,
		Closer:       closer,
		doer:         doer,
	}
}

// ListResources calls the ListResources API server method with
// the given service names.
func (c Client) ListResources(services []string) ([][]resource.Resource, error) {
	args, err := api.NewListResourcesArgs(services...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var apiResults api.ResourcesResults
	if err := c.FacadeCall("ListResources", &args, &apiResults); err != nil {
		return nil, errors.Trace(err)
	}

	if len(apiResults.Results) != len(services) {
		// We don't bother returning the results we *did* get since
		// something bad happened on the server.
		return nil, errors.Errorf("got invalid data from server (expected %d results, got %d)", len(services), len(apiResults.Results))
	}

	results := make([][]resource.Resource, len(services))
	for i := range services {
		apiResult := apiResults.Results[i]

		result, err := api.APIResult2Resources(apiResult)
		if err != nil {
			// TODO(ericsnow) Aggregate errors?
			return nil, errors.Trace(err)
		}
		results[i] = result
	}

	return results, nil
}

func (c Client) Upload(service, name string, reader io.ReadSeeker) error {
	// hash resource
	// upload resource

	args, err := api.NewUploadArgs(service, name)
	if err != nil {
		return errors.Trace(err)
	}

	req, err := http.NewRequest("PUT", "/"+resource.ComponentName, nil)
	if err != nil {
		return errors.Trace(err)
	}
	body, contentType, err := httpattachment.NewBody(reader, args, "juju-resource-"+service+"-"+name)
	if err != nil {
		return errors.Annotatef(err, "cannot create multipart body")
	}
	req.Header.Set("Content-Type", contentType)
	if err := c.doer.Do(req, body, nil); err != nil {
		return errors.Trace(err)
	}
	return nil
}
