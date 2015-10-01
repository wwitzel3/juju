// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"sync"

	"github.com/juju/errors"

	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/provider/common"
)

type baseProvider interface {
	// BootstrapEnv bootstraps a Juju environment.
	BootstrapEnv(environs.BootstrapContext, environs.BootstrapParams) (string, string, environs.BootstrapFinalizer, error)

	// DestroyEnv destroys the provided Juju environment.
	DestroyEnv() error
}

type environ struct {
	common.SupportsUnitPlacementPolicy

	name string
	uuid string
	raw  *rawProvider
	base baseProvider

	lock sync.Mutex
	ecfg *environConfig
}

type newRawProviderFunc func(*environConfig) (*rawProvider, error)

func newEnviron(cfg *config.Config, newRawProvider newRawProviderFunc) (*environ, error) {
	ecfg, err := newValidConfig(cfg, configDefaults)
	if err != nil {
		return nil, errors.Annotate(err, "invalid config")
	}

	// Connect and authenticate.
	raw, err := newRawProvider(ecfg)
	if err != nil {
		return nil, errors.Trace(err)
	}

	env, err := newEnvironRaw(ecfg, raw)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return env, nil
}

func newEnvironRaw(ecfg *environConfig, raw *rawProvider) (*environ, error) {
	uuid, ok := ecfg.UUID()
	if !ok {
		return nil, errors.New("UUID not set")
	}

	env := &environ{
		name: ecfg.Name(),
		uuid: uuid,
		ecfg: ecfg,
		raw:  raw,
	}
	env.base = common.DefaultProvider{Env: env}
	return env, nil
}

// Name returns the name of the environment.
func (env *environ) Name() string {
	return env.name
}

// Provider returns the environment provider that created this env.
func (*environ) Provider() environs.EnvironProvider {
	return providerInstance
}

// SetConfig updates the env's configuration.
func (env *environ) SetConfig(cfg *config.Config) error {
	env.lock.Lock()
	defer env.lock.Unlock()

	if env.ecfg == nil {
		return errors.New("cannot set config on uninitialized env")
	}

	if err := env.ecfg.update(cfg); err != nil {
		return errors.Annotate(err, "invalid config change")
	}
	return nil
}

// getSnapshot returns a copy of the environment. This is useful for
// ensuring the env you are using does not get changed by other code
// while you are using it.
func (env environ) getSnapshot() *environ {
	return &env
}

// Config returns the configuration data with which the env was created.
func (env *environ) Config() *config.Config {
	return env.getSnapshot().ecfg.Config
}

// Bootstrap creates a new instance, chosing the series and arch out of
// available tools. The series and arch are returned along with a func
// that must be called to finalize the bootstrap process by transferring
// the tools and installing the initial juju state server.
func (env *environ) Bootstrap(ctx environs.BootstrapContext, params environs.BootstrapParams) (arch, series string, _ environs.BootstrapFinalizer, _ error) {
	// TODO(ericsnow) Ensure currently not the root user
	// if remote is local host?

	// Using the Bootstrap func from provider/common should be fine.
	// Local provider does its own thing because it has to deal directly
	// with localhost rather than using SSH.
	return env.base.BootstrapEnv(ctx, params)
}

// Destroy shuts down all known machines and destroys the rest of the
// known environment.
func (env *environ) Destroy() error {
	ports, err := env.Ports()
	if err != nil {
		return errors.Trace(err)
	}

	if len(ports) > 0 {
		if err := env.ClosePorts(ports); err != nil {
			return errors.Trace(err)
		}
	}

	if err := env.base.DestroyEnv(); err != nil {
		return errors.Trace(err)
	}
	return nil
}