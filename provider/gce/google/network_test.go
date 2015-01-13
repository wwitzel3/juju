// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package google_test

import (
	"code.google.com/p/google-api-go-client/compute/v1"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/network"
	"github.com/juju/juju/provider/gce/google"
)

type networkSuite struct {
	google.BaseSuite
}

var _ = gc.Suite(&networkSuite{})

func (s *networkSuite) TestNetworkSpecPath(c *gc.C) {
	spec := google.NetworkSpec{
		Name: "spam",
	}
	path := spec.Path()

	c.Check(path, gc.Equals, "global/networks/spam")
}

func (s *networkSuite) TestNetworkSpecNewInterface(c *gc.C) {
	spec := google.NetworkSpec{
		Name: "spam",
	}
	netIF := google.NewNetInterface(spec, "eggs")

	c.Check(netIF, gc.DeepEquals, &compute.NetworkInterface{
		Network: "global/networks/spam",
		AccessConfigs: []*compute.AccessConfig{{
			Name: "eggs",
			Type: google.NetworkAccessOneToOneNAT,
		}},
	})
}

func (s *networkSuite) TestFirewallSpec(c *gc.C) {
	ports := network.NewPortSet(
		network.MustParsePortRange("80-81/tcp"),
		network.MustParsePortRange("8888/tcp"),
		network.MustParsePortRange("1234/udp"),
	)
	fw := google.FirewallSpec("spam", ports)

	allowed := []*compute.FirewallAllowed{{
		IPProtocol: "tcp",
		Ports:      []string{"80", "81", "8888"},
	}, {
		IPProtocol: "udp",
		Ports:      []string{"1234"},
	}}
	c.Check(fw, jc.DeepEquals, &compute.Firewall{
		Name:         "spam",
		TargetTags:   []string{"spam"},
		SourceRanges: []string{"0.0.0.0/0"},
		Allowed:      allowed,
	})
}