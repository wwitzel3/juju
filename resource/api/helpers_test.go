// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package api_test

import (
	"time"

	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/resource"
	"github.com/juju/juju/resource/api"
)

const fingerprint = "123456789012345678901234567890123456789012345678"

func newFingerprint(c *gc.C, data string) charmresource.Fingerprint {
	fp, err := charmresource.GenerateFingerprint([]byte(data))
	c.Assert(err, jc.ErrorIsNil)
	return fp
}

type helpersSuite struct {
	testing.IsolationSuite
}

var _ = gc.Suite(&helpersSuite{})

func (helpersSuite) TestResource2API(c *gc.C) {
	fp, err := charmresource.NewFingerprint([]byte(fingerprint))
	c.Assert(err, jc.ErrorIsNil)
	now := time.Now()
	res := resource.Resource{
		Resource: charmresource.Resource{
			Meta: charmresource.Meta{
				Name:    "spam",
				Type:    charmresource.TypeFile,
				Path:    "spam.tgz",
				Comment: "you need it",
			},
			Origin:      charmresource.OriginUpload,
			Revision:    1,
			Fingerprint: fp,
		},
		Username:  "a-user",
		Timestamp: now,
	}
	err = res.Validate()
	c.Assert(err, jc.ErrorIsNil)
	apiRes := api.Resource2API(res)

	c.Check(apiRes, jc.DeepEquals, api.Resource{
		CharmResource: api.CharmResource{
			Name:        "spam",
			Type:        "file",
			Path:        "spam.tgz",
			Comment:     "you need it",
			Origin:      "upload",
			Revision:    1,
			Fingerprint: []byte(fingerprint),
		},
		Username:  "a-user",
		Timestamp: now,
	})
}

func (helpersSuite) TestAPIResult2ResourcesOkay(c *gc.C) {
	fp, err := charmresource.NewFingerprint([]byte(fingerprint))
	c.Assert(err, jc.ErrorIsNil)
	now := time.Now()
	expected := resource.Resource{
		Resource: charmresource.Resource{
			Meta: charmresource.Meta{
				Name:    "spam",
				Type:    charmresource.TypeFile,
				Path:    "spam.tgz",
				Comment: "you need it",
			},
			Origin:      charmresource.OriginUpload,
			Revision:    1,
			Fingerprint: fp,
		},
		Username:  "a-user",
		Timestamp: now,
	}
	err = expected.Validate()
	c.Assert(err, jc.ErrorIsNil)
	apiRes := api.Resource{
		CharmResource: api.CharmResource{
			Name:        "spam",
			Type:        "file",
			Path:        "spam.tgz",
			Comment:     "you need it",
			Origin:      "upload",
			Revision:    1,
			Fingerprint: []byte(fingerprint),
		},
		Username:  "a-user",
		Timestamp: now,
	}

	resources, err := api.APIResult2Resources(api.ResourcesResult{
		Resources: []api.Resource{
			apiRes,
		},
	})
	c.Assert(err, jc.ErrorIsNil)

	c.Check(resources, jc.DeepEquals, []resource.Resource{
		expected,
	})
}

func (helpersSuite) TestAPIResult2ResourcesFailure(c *gc.C) {
	apiRes := api.Resource{
		CharmResource: api.CharmResource{
			Name:        "spam",
			Type:        "file",
			Path:        "spam.tgz",
			Origin:      "upload",
			Revision:    1,
			Fingerprint: []byte(fingerprint),
		},
	}
	failure := errors.New("<failure>")

	_, err := api.APIResult2Resources(api.ResourcesResult{
		ErrorResult: params.ErrorResult{
			Error: &params.Error{
				Message: failure.Error(),
			},
		},
		Resources: []api.Resource{
			apiRes,
		},
	})

	c.Check(err, gc.ErrorMatches, "<failure>")
	c.Check(errors.Cause(err), gc.Not(gc.Equals), failure)
}

func (helpersSuite) TestAPIResult2ResourcesNotFound(c *gc.C) {
	apiRes := api.Resource{
		CharmResource: api.CharmResource{
			Name:        "spam",
			Type:        "file",
			Path:        "spam.tgz",
			Origin:      "upload",
			Revision:    1,
			Fingerprint: []byte(fingerprint),
		},
	}

	_, err := api.APIResult2Resources(api.ResourcesResult{
		ErrorResult: params.ErrorResult{
			Error: &params.Error{
				Message: `service "a-service" not found`,
				Code:    params.CodeNotFound,
			},
		},
		Resources: []api.Resource{
			apiRes,
		},
	})

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (helpersSuite) TestAPI2Resource(c *gc.C) {
	now := time.Now()
	res, err := api.API2Resource(api.Resource{
		CharmResource: api.CharmResource{
			Name:        "spam",
			Type:        "file",
			Path:        "spam.tgz",
			Comment:     "you need it",
			Origin:      "upload",
			Revision:    1,
			Fingerprint: []byte(fingerprint),
		},
		Username:  "a-user",
		Timestamp: now,
	})
	c.Assert(err, jc.ErrorIsNil)

	fp, err := charmresource.NewFingerprint([]byte(fingerprint))
	c.Assert(err, jc.ErrorIsNil)
	expected := resource.Resource{
		Resource: charmresource.Resource{
			Meta: charmresource.Meta{
				Name:    "spam",
				Type:    charmresource.TypeFile,
				Path:    "spam.tgz",
				Comment: "you need it",
			},
			Origin:      charmresource.OriginUpload,
			Revision:    1,
			Fingerprint: fp,
		},
		Username:  "a-user",
		Timestamp: now,
	}
	err = expected.Validate()
	c.Assert(err, jc.ErrorIsNil)

	c.Check(res, jc.DeepEquals, expected)
}

func (helpersSuite) TestCharmResource2API(c *gc.C) {
	fp, err := charmresource.NewFingerprint([]byte(fingerprint))
	c.Assert(err, jc.ErrorIsNil)
	res := charmresource.Resource{
		Meta: charmresource.Meta{
			Name:    "spam",
			Type:    charmresource.TypeFile,
			Path:    "spam.tgz",
			Comment: "you need it",
		},
		Origin:      charmresource.OriginUpload,
		Revision:    1,
		Fingerprint: fp,
	}
	err = res.Validate()
	c.Assert(err, jc.ErrorIsNil)
	apiInfo := api.CharmResource2API(res)

	c.Check(apiInfo, jc.DeepEquals, api.CharmResource{
		Name:        "spam",
		Type:        "file",
		Path:        "spam.tgz",
		Comment:     "you need it",
		Origin:      "upload",
		Revision:    1,
		Fingerprint: []byte(fingerprint),
	})
}

func (helpersSuite) TestAPI2CharmResource(c *gc.C) {
	res, err := api.API2CharmResource(api.CharmResource{
		Name:        "spam",
		Type:        "file",
		Path:        "spam.tgz",
		Comment:     "you need it",
		Origin:      "upload",
		Revision:    1,
		Fingerprint: []byte(fingerprint),
	})
	c.Assert(err, jc.ErrorIsNil)

	fp, err := charmresource.NewFingerprint([]byte(fingerprint))
	c.Assert(err, jc.ErrorIsNil)
	expected := charmresource.Resource{
		Meta: charmresource.Meta{
			Name:    "spam",
			Type:    charmresource.TypeFile,
			Path:    "spam.tgz",
			Comment: "you need it",
		},
		Origin:      charmresource.OriginUpload,
		Revision:    1,
		Fingerprint: fp,
	}
	err = expected.Validate()
	c.Assert(err, jc.ErrorIsNil)

	c.Check(res, jc.DeepEquals, expected)
}
