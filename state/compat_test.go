// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"labix.org/v2/mgo/bson"
	"labix.org/v2/mgo/txn"
	gc "launchpad.net/gocheck"

	"github.com/juju/juju/network"
	"github.com/juju/juju/state"

	jujutesting "github.com/juju/juju/juju/testing"
)

// compatSuite contains backwards compatibility tests,
// for ensuring state operations behave correctly across
// schema changes.
type compatSuite struct {
	jujutesting.RepoSuite
}

var _ = gc.Suite(&compatSuite{})

func (s *compatSuite) TestEnvironAssertAlive(c *gc.C) {
	env, err := s.State.Environment()
	c.Assert(err, gc.IsNil)

	// 1.17+ has a "Life" field in environment documents.
	// We remove it here, to test 1.16 compatibility.
	ops := []txn.Op{{
		C:      env.Name(),
		Id:     env.UUID(),
		Update: bson.D{{"$unset", bson.D{{"life", nil}}}},
	}}
	err = state.RunTransaction(s.State, ops)
	c.Assert(err, gc.IsNil)

	// Now check the assertAliveOp and Destroy work as if
	// the environment is Alive.
	err = state.RunTransaction(s.State, []txn.Op{state.AssertAliveOp(env)})
	c.Assert(err, gc.IsNil)
	err = env.Destroy()
	c.Assert(err, gc.IsNil)
}

func (s *compatSuite) TestGetServiceWithoutNetworksIsOK(c *gc.C) {
	_, err := s.State.AddAdminUser("pass")
	c.Assert(err, gc.IsNil)
	charm := state.AddTestingCharm(c, s.State, "mysql")
	service, err := s.State.AddService("mysql", "user-admin", charm, nil)
	c.Assert(err, gc.IsNil)
	// In 1.17.7+ all services have associated document in the
	// requested networks collection. We remove it here to test
	// backwards compatibility.
	ops := []txn.Op{state.RemoveRequestedNetworksOp(s.State, state.GlobalKey(service))}
	err = state.RunTransaction(s.State, ops)
	c.Assert(err, gc.IsNil)

	// Now check the trying to fetch service's networks is OK.
	networks, err := service.Networks()
	c.Assert(err, gc.IsNil)
	c.Assert(networks, gc.HasLen, 0)
}

func (s *compatSuite) TestGetMachineWithoutRequestedNetworksIsOK(c *gc.C) {
	machine, err := s.State.EnvironmentDeployer.AddMachine("quantal", state.JobHostUnits)
	c.Assert(err, gc.IsNil)
	// In 1.17.7+ all machines have associated document in the
	// requested networks collection. We remove it here to test
	// backwards compatibility.
	ops := []txn.Op{state.RemoveRequestedNetworksOp(s.State, state.GlobalKey(machine))}
	err = state.RunTransaction(s.State, ops)
	c.Assert(err, gc.IsNil)

	// Now check the trying to fetch machine's networks is OK.
	networks, err := machine.RequestedNetworks()
	c.Assert(err, gc.IsNil)
	c.Assert(networks, gc.HasLen, 0)
}

// Check if ports stored on the unit are displayed.
func (s *compatSuite) TestShowUnitPorts(c *gc.C) {
	_, err := s.state.AddAdminUser("pass")
	c.Assert(err, gc.IsNil)
	charm := addCharm(c, s.state, "quantal", charmtesting.Charms.Dir("mysql"))
	service, err := s.state.AddService("mysql", "user-admin", charm, nil)
	c.Assert(err, gc.IsNil)
	unit, err := service.AddUnit()
	c.Assert(err, gc.IsNil)
	machine, err := s.state.AddMachine("quantal", JobHostUnits)
	c.Assert(err, gc.IsNil)
	c.Assert(unit.AssignToMachine(machine), gc.IsNil)

	// Add old-style ports to unit.
	port := network.Port{Protocol: "tcp", Number: 80}
	ops := []txn.Op{{
		C:      s.state.units.Name,
		Id:     unit.doc.Name,
		Assert: notDeadDoc,
		Update: bson.D{{"$addToSet", bson.D{{"ports", port}}}},
	}}
	err = s.state.runTransaction(ops)
	c.Assert(err, gc.IsNil)
	err = unit.Refresh()
	c.Assert(err, gc.IsNil)

	ports := unit.OpenedPorts()
	c.Assert(ports, gc.DeepEquals, []network.Port{{"tcp", 80}})
}

// Check if opening ports on a unit with ports stored in the unit doc works.
func (s *compatSuite) TestMigratePortsOnOpen(c *gc.C) {
	_, err := s.state.AddAdminUser("pass")
	c.Assert(err, gc.IsNil)
	charm := addCharm(c, s.state, "quantal", charmtesting.Charms.Dir("mysql"))
	service, err := s.state.AddService("mysql", "user-admin", charm, nil)
	c.Assert(err, gc.IsNil)
	unit, err := service.AddUnit()
	c.Assert(err, gc.IsNil)
	machine, err := s.state.AddMachine("quantal", JobHostUnits)
	c.Assert(err, gc.IsNil)
	c.Assert(unit.AssignToMachine(machine), gc.IsNil)

	// Add old-style ports to unit.
	port := network.Port{Protocol: "tcp", Number: 80}
	ops := []txn.Op{{
		C:      s.state.units.Name,
		Id:     unit.doc.Name,
		Assert: notDeadDoc,
		Update: bson.D{{"$addToSet", bson.D{{"ports", port}}}},
	}}
	err = s.state.runTransaction(ops)
	c.Assert(err, gc.IsNil)
	err = unit.Refresh()
	c.Assert(err, gc.IsNil)

	// Check if port conflicts are detected.
	err = unit.OpenPort("tcp", 80)
	c.Assert(err, gc.ErrorMatches, "cannot open ports 80-80/tcp for unit \"mysql/0\": cannot open ports 80-80/tcp on machine 0 due to conflict")

	err = unit.OpenPort("tcp", 8080)
	c.Assert(err, gc.IsNil)

	ports := unit.OpenedPorts()
	c.Assert(ports, gc.DeepEquals, []network.Port{{"tcp", 80}, {"tcp", 8080}})
}

// Check if closing ports on a unit with ports stored in the unit doc works.
func (s *compatSuite) TestMigratePortsOnClose(c *gc.C) {
	_, err := s.state.AddAdminUser("pass")
	c.Assert(err, gc.IsNil)
	charm := addCharm(c, s.state, "quantal", charmtesting.Charms.Dir("mysql"))
	service, err := s.state.AddService("mysql", "user-admin", charm, nil)
	c.Assert(err, gc.IsNil)
	unit, err := service.AddUnit()
	c.Assert(err, gc.IsNil)
	machine, err := s.state.AddMachine("quantal", JobHostUnits)
	c.Assert(err, gc.IsNil)
	c.Assert(unit.AssignToMachine(machine), gc.IsNil)

	// Add old-style ports to unit.
	port := network.Port{Protocol: "tcp", Number: 80}
	ops := []txn.Op{{
		C:      s.state.units.Name,
		Id:     unit.doc.Name,
		Assert: notDeadDoc,
		Update: bson.D{{"$addToSet", bson.D{{"ports", port}}}},
	}}
	err = s.state.runTransaction(ops)
	c.Assert(err, gc.IsNil)
	err = unit.Refresh()
	c.Assert(err, gc.IsNil)

	// Check if closing an unopened port causes error
	err = unit.ClosePort("tcp", 8080)
	c.Assert(err, gc.ErrorMatches, "cannot close ports 8080-8080/tcp for unit \"mysql/0\": no match found for port range: 8080-8080/tcp")

	err = unit.ClosePort("tcp", 80)
	c.Assert(err, gc.IsNil)

	ports := unit.OpenedPorts()
	c.Assert(ports, gc.DeepEquals, []network.Port{})
}
