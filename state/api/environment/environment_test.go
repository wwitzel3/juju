// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environment_test

import (
	gc "launchpad.net/gocheck"

	jujutesting "github.com/juju/juju/juju/testing"
	"github.com/juju/juju/state/api"
	"github.com/juju/juju/state/api/environment"
	apitesting "github.com/juju/juju/state/api/testing"
)

type environmentSuite struct {
	jujutesting.JujuConnSuite
	*apitesting.EnvironWatcherTests

	st          *api.State
	environment *environment.Facade
}

var _ = gc.Suite(&environmentSuite{})

func (s *environmentSuite) SetUpTest(c *gc.C) {
	s.JujuConnSuite.SetUpTest(c)

	s.st, _ = s.OpenAPIAsNewMachine(c)

	s.environment = s.st.Environment()
	c.Assert(s.environment, gc.NotNil)

	s.EnvironWatcherTests = apitesting.NewEnvironWatcherTests(
		s.environment, s.BackingState, apitesting.NoSecrets)
}

func (s *environmentSuite) TestGetCapabilities(c *gc.C) {
	result, err := s.environment.GetCapabilities()
	c.Assert(err, gc.IsNil)
	c.Assert(result, gc.NotNil)

	c.Assert(result.SupportedArchitectures, gc.DeepEquals, []string{"amd64", "i386", "ppc64"})
	c.Assert(result.AvailabilityZones[0], gc.DeepEquals, map[string]interface{}{"available": false, "name": "zone_0"})

	instanceType := result.InstanceTypes[0]
	c.Assert(instanceType["name"], gc.Equals, "test-name")
}
