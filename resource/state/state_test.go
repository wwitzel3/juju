// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"github.com/juju/errors"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/resource/state"
)

var _ = gc.Suite(&StateSuite{})

type StateSuite struct {
	testing.IsolationSuite

	stub    *testing.Stub
	raw     *stubRawState
	persist *stubPersistence
	storage *stubStorage
}

func (s *StateSuite) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	s.stub = &testing.Stub{}
	s.raw = &stubRawState{stub: s.stub}
	s.persist = &stubPersistence{stub: s.stub}
	s.storage = &stubStorage{stub: s.stub}
	s.raw.ReturnPersistence = s.persist
	s.raw.ReturnStorage = s.storage
}

func (s *StateSuite) TestNewStateOkay(c *gc.C) {
	_, err := state.NewState(s.raw)
	c.Assert(err, jc.ErrorIsNil)

	s.stub.CheckCallNames(c, "Persistence", "Storage")
}

func (s *StateSuite) TestNewStateFailure(c *gc.C) {
	failure := errors.New("<failure>")
	s.stub.SetErrors(failure)

	_, err := state.NewState(s.raw)

	c.Check(errors.Cause(err), gc.Equals, failure)
}
