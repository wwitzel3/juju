// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common_test

import (
	"fmt"

	"github.com/juju/juju/instance"

	"github.com/juju/juju/provider/common"
	coretesting "github.com/juju/juju/testing"

	gc "launchpad.net/gocheck"
)

type EnvironCapabilitiesSuite struct {
	coretesting.FakeJujuHomeSuite
	env mockZonedEnviron
}

var _ = gc.Suite(&EnvironCapabilitiesSuite{})

func (s *EnvironCapabilitiesSuite) SetUpTest(c *gc.C) {
	s.FakeJujuHomeSuite.SetUpSuite(c)

	allInstances := make([]instance.Instance, 3)
	for i := range allInstances {
		allInstances[i] = &mockInstance{id: fmt.Sprintf("inst%d", i)}
	}
	s.env.allInstances = func() ([]instance.Instance, error) {
		return allInstances, nil
	}

	availabilityZones := make([]common.AvailabilityZone, 3)
	for i := range availabilityZones {
		availabilityZones[i] = &mockAvailabilityZone{
			name:      fmt.Sprintf("az%d", i),
			available: i > 0,
		}
	}
	s.env.availabilityZones = func() ([]common.AvailabilityZone, error) {
		return availabilityZones, nil
	}
}

func (s *EnvironCapabilitiesSuite) TestNewEnvironCapabilities(c *gc.C) {
	environCapabilities, err := common.NewEnvironCapabilities(&s.env)
	c.Assert(err, gc.Equals, nil)
	c.Assert(environCapabilities, gc.Not(gc.Equals), nil)

	data := environCapabilities.Result()
	c.Assert(data, gc.NotNil)
}
