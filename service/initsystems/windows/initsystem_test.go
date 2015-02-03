// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package windows_test

import (
	"fmt"

	"github.com/juju/errors"
	//"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils"
	"github.com/juju/utils/fs"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/service/initsystems"
	"github.com/juju/juju/service/initsystems/windows"
	coretesting "github.com/juju/juju/testing"
)

const confStr = `{
 "description": "juju agent for %s",
 "startexec": "jujud.exe %s"
}`

type initSystemSuite struct {
	coretesting.BaseSuite

	initDir string
	files   *fs.FakeOps
	cmd     *initsystems.FakeShell
	init    initsystems.InitSystem
	conf    initsystems.Conf
	confStr string
}

var _ = gc.Suite(&initSystemSuite{})

func (s *initSystemSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)

	s.files = fs.NewFakeOps()
	s.cmd = &initsystems.FakeShell{}
	s.init = windows.NewWindows(s.files, s.cmd)
	s.conf = initsystems.Conf{
		Desc: "juju agent for jujud-machine-0",
		Cmd:  "jujud.exe machine-0",
	}
	s.confStr = s.newConfStr("jujud-machine-0")

	s.PatchValue(&initsystems.RetryAttempts, utils.AttemptStrategy{})
}

func (s *initSystemSuite) newConfStr(name string) string {
	tag := name[len("jujud-"):]
	return fmt.Sprintf(confStr, name, tag)
}

func (s *initSystemSuite) setStatus(name, status string) {
	switch status {
	case initsystems.StatusRunning:
		s.cmd.Out = []byte("Running\n")
	case initsystems.StatusEnabled:
		s.cmd.Out = []byte("Stopped\n")
	case "":
		s.cmd.SetErrors(errors.New("service " + name + " not found"))
	}
}

func (s *initSystemSuite) TestInitSystemName(c *gc.C) {
	name := s.init.Name()

	c.Check(name, gc.Equals, "windows")
}

func (s *initSystemSuite) TestInitSystemList(c *gc.C) {
	s.cmd.Out = []byte("" +
		"jujud-machine-0 " +
		"something-else " +
		"jujud-unit-wordpress-0")

	names, err := s.init.List()
	c.Assert(err, jc.ErrorIsNil)

	c.Check(names, jc.SameContents, []string{
		"jujud-machine-0",
		"something-else",
		"jujud-unit-wordpress-0",
	})
}

func (s *initSystemSuite) TestInitSystemListLimited(c *gc.C) {
	s.cmd.Out = []byte("" +
		"jujud-machine-0 " +
		"something-else " +
		"jujud-unit-wordpress-0")

	names, err := s.init.List("jujud-machine-0")
	c.Assert(err, jc.ErrorIsNil)

	c.Check(names, jc.SameContents, []string{"jujud-machine-0"})
}

func (s *initSystemSuite) TestInitSystemListLimitedEmpty(c *gc.C) {
	names, err := s.init.List("jujud-machine-0")
	c.Assert(err, jc.ErrorIsNil)

	c.Check(names, jc.SameContents, []string{})
}

func (s *initSystemSuite) TestInitSystemStart(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	err := s.init.Start(name)
	c.Assert(err, jc.ErrorIsNil)

	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemStartAlreadyRunning(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusRunning)

	err := s.init.Start(name)

	c.Check(err, jc.Satisfies, errors.IsAlreadyExists)
}

func (s *initSystemSuite) TestInitSystemStartNotEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	err := s.init.Start(name)

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemStop(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusRunning)

	err := s.init.Stop(name)
	c.Assert(err, jc.ErrorIsNil)

	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemStopNotRunning(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	err := s.init.Stop(name)

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemStopNotEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	err := s.init.Stop("jujud-unit-wordpress-0")

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemEnable(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")
	s.files.Returns.Data = []byte(s.newConfStr(name))

	filename := "/var/lib/juju/init/jujud-machine-0"
	err := s.init.Enable(name, filename)
	c.Assert(err, jc.ErrorIsNil)

	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemEnableAlreadyEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	filename := "/var/lib/juju/init/jujud-machine-0"
	err := s.init.Enable(name, filename)

	c.Check(err, jc.Satisfies, errors.IsAlreadyExists)
}

func (s *initSystemSuite) TestInitSystemDisable(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	err := s.init.Disable(name)
	c.Assert(err, jc.ErrorIsNil)

	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemDisableNotEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	err := s.init.Disable(name)

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemIsEnabledTrue(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	enabled, err := s.init.IsEnabled(name)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(enabled, jc.IsTrue)
	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemIsEnabledFalse(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	enabled, err := s.init.IsEnabled(name)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(enabled, jc.IsFalse)
}

func (s *initSystemSuite) TestInitSystemInfoRunning(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusRunning)

	info, err := s.init.Info(name)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(info, jc.DeepEquals, &initsystems.ServiceInfo{
		Name:   name,
		Status: initsystems.StatusRunning,
	})
	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemInfoNotRunning(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	info, err := s.init.Info(name)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(info, jc.DeepEquals, &initsystems.ServiceInfo{
		Name:   name,
		Status: initsystems.StatusStopped,
	})
	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemInfoNotEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	_, err := s.init.Info(name)

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemConf(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, initsystems.StatusEnabled)

	conf, err := s.init.Conf(name)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(conf, jc.DeepEquals, &initsystems.Conf{
		Desc: `juju agent for jujud-unit-wordpress-0`,
		Cmd:  "jujud.exe unit-wordpress-0",
	})
	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemConfNotEnabled(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	s.setStatus(name, "")

	_, err := s.init.Conf(name)

	c.Check(err, jc.Satisfies, errors.IsNotFound)
}

func (s *initSystemSuite) TestInitSystemValidate(c *gc.C) {
	err := s.init.Validate("jujud-machine-0", s.conf)
	c.Assert(err, jc.ErrorIsNil)

	// TODO(ericsnow) Check underlying calls.
}

func (s *initSystemSuite) TestInitSystemValidateInvalid(c *gc.C) {
	s.conf.Cmd = ""

	err := s.init.Validate("jujud-machine-0", s.conf)

	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *initSystemSuite) TestInitSystemValidateUnsupportedEnv(c *gc.C) {
	s.conf.Env = map[string]string{
		"x": "y",
	}

	err := s.init.Validate("jujud-machine-0", s.conf)

	expected := initsystems.NewUnsupportedField("Env")
	c.Check(errors.Cause(err), gc.FitsTypeOf, expected)
}

func (s *initSystemSuite) TestInitSystemValidateUnsupportedLimit(c *gc.C) {
	s.conf.Limit = map[string]string{
		"x": "y",
	}

	err := s.init.Validate("jujud-machine-0", s.conf)

	expected := initsystems.NewUnsupportedField("Limit")
	c.Check(errors.Cause(err), gc.FitsTypeOf, expected)
}

func (s *initSystemSuite) TestInitSystemValidateUnsupportedOut(c *gc.C) {
	s.conf.Out = "/var/log/juju/machine-0.log"

	err := s.init.Validate("jujud-machine-0", s.conf)

	expected := initsystems.NewUnsupportedField("Out")
	c.Check(errors.Cause(err), gc.FitsTypeOf, expected)
}

func (s *initSystemSuite) TestInitSystemSerialize(c *gc.C) {
	data, err := s.init.Serialize("jujud-machine-0", s.conf)
	c.Assert(err, jc.ErrorIsNil)

	c.Check(string(data), gc.Equals, s.confStr)
}

func (s *initSystemSuite) TestInitSystemSerializeUnsupported(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	conf := initsystems.Conf{
		Desc: "juju agent for " + name,
		Cmd:  "jujud.exe unit-wordpress-0",
		Out:  "/var/log/juju/unit-wordpress-0",
	}
	_, err := s.init.Serialize("jujud-machine-0", conf)

	expected := initsystems.NewUnsupportedField("Out")
	c.Check(errors.Cause(err), gc.FitsTypeOf, expected)
}

func (s *initSystemSuite) TestInitSystemDeserialize(c *gc.C) {
	name := "jujud-unit-wordpress-0"
	data := s.newConfStr(name)
	conf, err := s.init.Deserialize([]byte(data))
	c.Assert(err, jc.ErrorIsNil)

	c.Check(conf, jc.DeepEquals, &initsystems.Conf{
		Desc: "juju agent for jujud-unit-wordpress-0",
		Cmd:  "jujud.exe unit-wordpress-0",
	})
}

func (s *initSystemSuite) TestInitSystemDeserializeUnsupported(c *gc.C) {
	data := `{
 "description": "juju agent for machine-0",
 "startexec": "jujud.exe machine-0",
 "out": "/var/log/juju/machine-0.log"
}`
	_, err := s.init.Deserialize([]byte(data))

	expected := initsystems.NewUnsupportedField("Out")
	c.Check(errors.Cause(err), gc.FitsTypeOf, expected)
}
