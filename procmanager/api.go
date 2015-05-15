// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

const ProcManager = "ProcManager"

type CharmCmds interface {
	// $ launch -> API.CharmHookEnv.Launch
	// valid arg parse
	// exec plugin with validated arguments
	//     -> plugin will use PluginResources to verify storage and networking
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
	RegisterProcess(info ProcessInfo) string
}

type Destroy interface {
	Info(UUID string) (ProcessInfo, error)
	UnregisterProcess(UUID string) error
}

type Plugin interface {
	Launch(image, desc, storage, networking, args string) (ProcessDetails, error)
	Destroy(UUID string) error
}

// processDetails holds information about the process that only the plugin
// can determine.
type ProcessDetails struct {
	// uniqueID is provided by the plugin as a guaranteed way
	// to identify the process.
	UniqueID string
	// status of the process
	Status string
}

// processInfo holds information about a process that Juju needs.
type ProcessInfo struct {
	Image  string
	Args   string
	Desc   string
	Plugin string

	// TODO(wwitzel3) determine specific details for storage and networking
	Storage    string
	Networking string

	Details ProcessDetails
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
	Networking(networking string) (networkInfo, error)
}
