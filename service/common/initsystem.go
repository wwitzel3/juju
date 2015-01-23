// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

// InitSystem represents the functionality provided by an init system.
// It encompasses all init services on the host, rather than just juju-
// managed ones.
type InitSystem interface {
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
	Info(name string) (*ServiceInfo, error)

	// Conf composes a Conf for the named service and returns it.
	// If the service is not enabled then errors.NotFound is returned.
	Conf(name string) (*Conf, error)

	// Serialize converts the provided Conf into the file format
	// recognized by the init system.
	Serialize(name string, conf *Conf) ([]byte, error)

	// Deserialize converts the provided data into a Conf according to
	// the init system's conf file format. If the data does not
	// correspond to that file format then an error is returned.
	Deserialize(data []byte) (*Conf, error)
}