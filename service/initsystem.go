// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"io/ioutil"
	"strings"

	"github.com/juju/errors"

	"github.com/juju/juju/service/common"
	"github.com/juju/juju/service/upstart"
	"github.com/juju/juju/service/windows"
	"github.com/juju/juju/version"
)

// These are the names of the juju-compatible init systems.
const (
	InitSystemWindows = "windows"
	InitSystemUpstart = "upstart"
	//InitSystemSystemd = "systemd"
)

var (
	linuxInitNames = map[string]string{
		"/sbin/init": InitSystemUpstart,
		//"/sbin/systemd": InitSystemSystemd,
	}
)

// InitSystem represents the functionality provided by an init system.
// It encompasses all init services on the host, rather than just juju-
// managed ones.
type InitSystem interface {
	// Name returns the init system's name.
	Name() string

	// List gathers the names of all enabled services in the init system
	// and returns them. If any names are passed as arguments then the
	// result will be limited to those names. Otherwise all known
	// service names are returned.
	List(include ...string) ([]string, error)

	// Start causes the named service to be started. If it is already
	// started then errors.AlreadyExists is returned. If the service has
	// not been enabled then errors.NotFound is returned.
	Start(name string) error

	// Stop causes the named service to be stopped. If it is already
	// stopped then errors.NotFound is returned. If the service has
	// not been enabled then errors.NotFound is returned.
	Stop(name string) error

	// Enable adds a new service to the init system with the given name.
	// The conf file at the provided filename is used for the new
	// service. If a service with that name is already enabled then
	// errors.AlreadyExists is returned.
	Enable(name, filename string) error

	// Disable removes the named service from the init system. If it is
	// not already enabled then errors.NotFound is returned.
	Disable(name string) error

	// IsEnabled determines whether or not the named service is enabled.
	IsEnabled(name string) (bool, error)

	// Info gathers information about the named service and returns it.
	// If the service is not enabled then errors.NotFound is returned.
	Info(name string) (*common.ServiceInfo, error)

	// Conf composes a Conf for the named service and returns it.
	// If the service is not enabled then errors.NotFound is returned.
	Conf(name string) (*common.Conf, error)

	// Validate checks the provided service name and conf to ensure
	// that they are compatible with the init system. If a particular
	// conf field is not supported by the init system then
	// errors.NotSupported is returned (see common.Conf). Otherwise
	// any other invalid results in an errors.NotValid error.
	Validate(name string, conf common.Conf) error

	// Serialize converts the provided Conf into the file format
	// recognized by the init system.
	Serialize(name string, conf common.Conf) ([]byte, error)

	// TODO(ericsnow) Pass name in or return it in Deserialize?

	// Deserialize converts the provided data into a Conf according to
	// the init system's conf file format. If the data does not
	// correspond to that file format then an error is returned.
	Deserialize(data []byte) (*common.Conf, error)
}

func newInitSystem(name string) InitSystem {
	switch name {
	case InitSystemWindows:
		return windows.NewInitSystem(name)
	case InitSystemUpstart:
		return upstart.NewInitSystem(name)
	}
	return nil
}

// discoverInitSystem determines which init system is running and
// returns its name.
func discoverInitSystem() (string, error) {
	if version.Current.OS == version.Windows {
		return InitSystemWindows, nil
	}

	executable, err := findInitExecutable()
	if err != nil {
		return "", errors.Annotate(err, "while finding init exe")
	}

	name, ok := linuxInitNames[executable]
	if !ok {
		return "", errors.New("unrecognized init system")
	}

	return name, nil
}

var findInitExecutable = func() (string, error) {
	data, err := ioutil.ReadFile("/proc/1/cmdline")
	if err != nil {
		return "", errors.Trace(err)
	}
	return strings.Fields(string(data))[0], nil
}
