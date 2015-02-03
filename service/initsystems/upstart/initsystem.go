// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upstart

import (
	"os"
	"path"
	"regexp"

	"github.com/juju/errors"

	"github.com/juju/juju/service/initsystems"
)

// Vars for patching in tests.
var (
	// ConfDir holds the default init directory name.
	ConfDir = "/etc/init"
)

var (
	upstartServicesRe = regexp.MustCompile("^([a-zA-Z0-9-_:]+)\\.conf$")
	upstartStartedRE  = regexp.MustCompile(`^.* start/running, process (\d+)\n$`)
)

type upstart struct {
	name    string
	initDir string
	fops    fileOperations
	cmd     cmdRunner
}

// NewInitSystem returns a new value that implements
// initsystems.InitSystem for upstart.
func NewInitSystem(name string) initsystems.InitSystem {
	return &upstart{
		name:    name,
		initDir: ConfDir,
		fops:    newFileOperations(),
		cmd:     newCmdRunner(),
	}
}

// confPath returns the path to the service's configuration file.
func (is upstart) confPath(name string) string {
	return path.Join(is.initDir, name+".conf")
}

// Name implements initsystems.InitSystem.
func (is upstart) Name() string {
	return is.name
}

// List implements initsystems.InitSystem.
func (is *upstart) List(include ...string) ([]string, error) {
	// TODO(ericsnow) We should be able to use initctl to do this.
	var services []string
	fis, err := is.fops.ListDir(is.initDir)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		groups := upstartServicesRe.FindStringSubmatch(fi.Name())
		if len(groups) > 0 {
			services = append(services, groups[1])
		}
	}

	return initsystems.FilterNames(services, include), nil
}

// Start implements initsystems.InitSystem.
func (is *upstart) Start(name string) error {
	if err := initsystems.EnsureEnabled(name, is); err != nil {
		return errors.Trace(err)
	}

	if is.isRunning(name) {
		return errors.AlreadyExistsf("service %q", name)
	}

	// On slower disks, upstart may take a short time to realise
	// that there is a service there.
	var err error
	for attempt := initsystems.RetryAttempts.Start(); attempt.Next(); {
		if err = is.start(name); err == nil {
			break
		}
	}
	return errors.Trace(err)
}

func (is *upstart) start(name string) error {
	_, err := is.cmd.RunCommand("start", "--system", name)
	if err != nil {
		// Double check to see if we were started before our command ran.
		if is.isRunning(name) {
			return nil
		}
		return errors.Trace(err)
	}
	return nil
}

// Stop implements initsystems.InitSystem.
func (is *upstart) Stop(name string) error {
	if err := initsystems.EnsureEnabled(name, is); err != nil {
		return errors.Trace(err)
	}

	if !is.isRunning(name) {
		return errors.NotFoundf("service %q", name)
	}

	_, err := is.cmd.RunCommand("stop", "--system", name)
	return errors.Trace(err)
}

// Enable implements initsystems.InitSystem.
func (is *upstart) Enable(name, filename string) error {
	// TODO(ericsnow) Deserialize and validate?

	enabled, err := is.IsEnabled(name)
	if err != nil {
		return errors.Trace(err)
	}
	if enabled {
		return errors.AlreadyExistsf("service %q", name)
	}

	err = is.fops.Symlink(filename, is.confPath(name))
	return errors.Trace(err)
}

// Disable implements initsystems.InitSystem.
func (is *upstart) Disable(name string) error {
	if err := initsystems.EnsureEnabled(name, is); err != nil {
		return errors.Trace(err)
	}

	enabled, err := is.IsEnabled(name)
	if err != nil {
		return errors.Trace(err)
	}
	if enabled {
		return nil
	}

	return os.Remove(is.confPath(name))
}

// TODO(ericsnow) Allow verifying against a file.

// IsEnabled implements initsystems.InitSystem.
func (is *upstart) IsEnabled(name string) (bool, error) {
	// TODO(ericsnow) In the general case, relying on the conf file
	// may not be the safest route. Perhaps we should use initctl?
	exists, err := is.fops.Exists(is.confPath(name))
	if err != nil {
		return false, errors.Trace(err)
	}
	return exists, nil
}

// Info implements initsystems.InitSystem.
func (is *upstart) Info(name string) (*initsystems.ServiceInfo, error) {
	if err := initsystems.EnsureEnabled(name, is); err != nil {
		return nil, errors.Trace(err)
	}

	conf, err := is.Conf(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	status := initsystems.StatusStopped
	if is.isRunning(name) {
		status = initsystems.StatusRunning
	}

	info := &initsystems.ServiceInfo{
		Name:        name,
		Description: conf.Desc,
		Status:      status,
	}
	return info, nil
}

func (is *upstart) isRunning(name string) bool {
	out, err := is.cmd.RunCommand("status", "--system", name)
	if err != nil {
		// TODO(ericsnow) Are we really okay ignoring the error?
		return false
	}
	return upstartStartedRE.Match(out)
}

// Conf implements initsystems.InitSystem.
func (is *upstart) Conf(name string) (*initsystems.Conf, error) {
	data, err := is.fops.ReadFile(is.confPath(name))
	if os.IsNotExist(err) {
		return nil, errors.NotFoundf("service %q", name)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}

	conf, err := is.Deserialize(data)
	return conf, errors.Trace(err)
}

// Validate implements initsystems.InitSystem.
func (is *upstart) Validate(name string, conf initsystems.Conf) error {
	err := Validate(name, conf)
	return errors.Trace(err)
}

// Serialize implements initsystems.InitSystem.
func (upstart) Serialize(name string, conf initsystems.Conf) ([]byte, error) {
	data, err := Serialize(name, conf)
	return data, errors.Trace(err)
}

// Deserialize implements initsystems.InitSystem.
func (is *upstart) Deserialize(data []byte) (*initsystems.Conf, error) {
	conf, err := Deserialize(data)
	return conf, errors.Trace(err)
}
