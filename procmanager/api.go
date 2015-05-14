// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

const ProcManagerAPI = "ProcManager"

type CharmCmds interface {
	// $ launch -> API.CharmHookEnv.Launch
	// valid arg parse
	// verify storage and networking with state
	// exec plugin with validated arguments
	// handle any plugin errors
	// convert unique identifier in to UUID for process
	// register UUID and process information with state
	Launch(plugin, image, desc, storage, networking, args string) (string, error) // UUID, error

	// $ destroy -> API.CharmHookEnv.Destroy
	// valid arg parse
	// verify UUID
	// exec plugin
	// handle/surface errors
	// unregister UUID with state
	Destroy(UUID string) error
}

type Launch interface {
	Verify(storage, networking string) error
	RegisterProcess(info processInfo) UUID
}

type Destroy interface {
	Info(UUID string) (processInfo, error)
	UnregisterProcess(UUID string) error
}

type Plugin interface {
	Launch(image, desc, storage, networking, args string) (processDetails, error)
	Destroy(UUID string) error
}

// processDetails holds information about the process that only the plugin
// can determine.
type processDetails struct {
	// uniqueID is provided by the plugin as a guaranteed way
	// to identify the process.
	uniqueID string
	// status of the process
	status string
}

// processInfo holds information about a process that Juju needs.
type processInfo struct {
	image  string
	args   string
	desc   string
	plugin string

	// TODO(wwitzel3) determine specific details for storage and networking
	storage    string
	networking string

	details processDetails
}

// TODO(wwitzel3) determine storageInfo based on spec/plugin needs
type storageInfo struct{}

// TODO(wwitzel3) determine networkInfo based on spec/plugin needs
type networkInfo struct{}

// PluginResource exposes an API for plugins to call
// to validate and allocate Juju resources
type PluginResource interface {
	// $ storage-info -> API.PluginResource.Storage
	Storage(storageID string) (storageInfo, error)
	// $ network-info -> API.PluginResource.Networking
	Networking(networking string) (networkingInfo, error)
}

// OLD STUFF
//
func init() {
	common.RegisterStandardFacade("ProcManager", 1, NewProcManagerAPI)
}

// ProcManagerAPI implements the API used by the procmanager worker.
type ProcManagerAPI struct {
	*common.EnvironWatcher

	st             *state.State
	resources      *common.Resources
	authorizer     common.Authorizer
	StateAddresser *common.StateAddresser
	canModify      bool
}

// NewProcManagerAPI creates a new instance of the ProcManager API.
func NewProcManagerAPI(st *state.State, resources *common.Resources, authorizer common.Authorizer) (*ProcManagerAPI, error) {
	if !authorizer.AuthMachineAgent() && !authorizer.AuthUnitAgent() {
		return nil, common.ErrPerm
	}
	return &ProcManagerAPI{
		EnvironWatcher: common.NewEnvironWatcher(st, resources, authorizer),
		st:             st,
		authorizer:     authorizer,
		resources:      resources,
		canModify:      authorizer.AuthEnvironManager(),
		StateAddresser: common.NewStateAddresser(st),
	}, nil
}
