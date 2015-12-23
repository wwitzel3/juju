// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package all

import (
	"io"
	"os"

	"github.com/juju/errors"
	"gopkg.in/juju/charm.v6-unstable"
	charmresource "gopkg.in/juju/charm.v6-unstable/resource"
	"gopkg.in/juju/charmrepo.v2-unstable"

	"github.com/juju/juju/api"
	"github.com/juju/juju/api/base"
	"github.com/juju/juju/apiserver/common"
	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/cmd/juju/commands"
	"github.com/juju/juju/resource"
	"github.com/juju/juju/resource/api/client"
	"github.com/juju/juju/resource/api/server"
	"github.com/juju/juju/resource/cmd"
	"github.com/juju/juju/resource/persistence"
	"github.com/juju/juju/resource/state"
	corestate "github.com/juju/juju/state"
)

// resources exposes the registration methods needed
// for the top-level component machinery.
type resources struct{}

// RegisterForServer is the top-level registration method
// for the component in a jujud context.
func (r resources) registerForServer() error {
	r.registerState()
	r.registerPublicFacade()
	return nil
}

// RegisterForClient is the top-level registration method
// for the component in a "juju" command context.
func (r resources) registerForClient() error {
	r.registerPublicCommands()
	return nil
}

// registerPublicFacade adds the resources public API facade
// to the API server.
func (r resources) registerPublicFacade() {
	if !markRegistered(resource.ComponentName, "public-facade") {
		return
	}

	common.RegisterStandardFacade(
		resource.ComponentName,
		server.Version,
		r.newPublicFacade,
	)
}

// newPublicFacade is passed into common.RegisterStandardFacade
// in registerPublicFacade.
func (resources) newPublicFacade(st *corestate.State, _ *common.Resources, authorizer common.Authorizer) (*server.Facade, error) {
	if !authorizer.AuthClient() {
		return nil, common.ErrPerm
	}

	rst, err := st.Resources()
	//rst, err := state.NewState(&resourceState{raw: st})
	if err != nil {
		return nil, errors.Trace(err)
	}

	return server.NewFacade(rst), nil
}

// resourcesApiClient adds a Close() method to the resources public API client.
type resourcesAPIClient struct {
	*client.Client
	closeConnFunc func() error
}

// Close implements io.Closer.
func (client resourcesAPIClient) Close() error {
	return client.closeConnFunc()
}

// registerState registers the state functionality for resources.
func (resources) registerState() {
	if !markRegistered(resource.ComponentName, "state") {
		return
	}

	newResources := func(persist corestate.Persistence) (corestate.Resources, error) {
		st, err := state.NewState(&resourceState{persist: persist})
		if err != nil {
			return nil, errors.Trace(err)
		}
		return st, nil
	}

	corestate.SetResourcesComponent(newResources)
}

// resourceState is a wrapper around state.State that supports the needs
// of resources.
type resourceState struct {
	persist corestate.Persistence
}

// Persistence implements resource/state.RawState.
func (st resourceState) Persistence() (state.Persistence, error) {
	return persistence.NewPersistence(st.persist), nil
}

// Storage implements resource/state.RawState.
func (st resourceState) Storage() (state.Storage, error) {
	return st.persist.NewStorage(), nil
}

// registerPublicCommands adds the resources-related commands
// to the "juju" supercommand.
func (r resources) registerPublicCommands() {
	if !markRegistered(resource.ComponentName, "public-commands") {
		return
	}

	newShowAPIClient := func(command *cmd.ShowCommand) (cmd.CharmResourceLister, error) {
		//return newCharmstore()
		store, err := newCharmstore()
		if err != nil {
			return nil, errors.Trace(err)
		}
		return &charmstore{store}, nil
	}
	commands.RegisterEnvCommand(func() envcmd.EnvironCommand {
		return cmd.NewShowCommand(newShowAPIClient)
	})

	commands.RegisterEnvCommand(func() envcmd.EnvironCommand {
		return cmd.NewUploadCommand(cmd.UploadDeps{
			NewClient: func(c *cmd.UploadCommand) (cmd.UploadClient, error) {
				return r.newClient(c)
			},
			OpenResource: func(s string) (io.ReadCloser, error) {
				return os.Open(s)
			},
		})

	})
}

func newCharmstore() (charmrepo.Interface, error) {
	// Also see apiserver/service/charmstore.go.
	var args charmrepo.NewCharmStoreParams
	store := charmrepo.NewCharmStore(args)
	return store, nil
}

// TODO(ericsnow) Get rid of charmstore one charmrepo.Interface grows the methods.

type charmstore struct {
	charmrepo.Interface
}

func (charmstore) ListResources(charmURLs []charm.URL) ([][]charmresource.Resource, error) {
	// TODO(ericsnow) finish!
	return nil, errors.Errorf("not implemented")
}

func (charmstore) Close() error {
	return nil
}

type apicommand interface {
	NewAPIRoot() (api.Connection, error)
}

func (resources) newClient(command apicommand) (*client.Client, error) {
	apiCaller, err := command.NewAPIRoot()
	if err != nil {
		return nil, errors.Trace(err)
	}
	caller := base.NewFacadeCallerForVersion(apiCaller, resource.ComponentName, server.Version)
	doer, err := apiCaller.HTTPClient()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return client.NewClient(caller, doer, apiCaller), nil
}
