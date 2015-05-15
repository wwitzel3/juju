// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

const ProcManager = "ProcManager"

// Plugin represents the functionality used by Juju of a process
// management plugin.
type Plugin interface {
	// TODO(wwitzel3) determine specific details for storage and networking
	// Launch launches a new process with the given info.
	Launch(image, desc, storage, networking, args string) (ProcessDetails, error)
	// Destroy destroys an existing process.
	Destroy(id string) error
	// Info retrieves information about an existing process, including status.
	Info(id string) (ProcessDetails, error)
}

// ProcessDetails holds information about an existing process as provided by
// the plugin.
type ProcessDetails struct {
	// UniqueID is provided by the plugin as a guaranteed way
	// to identify the process to the plugin.
	UniqueID string
	// Status is the status of the process as reported by the plugin.
	Status string
}

// ProcessInfo holds information about a process that Juju needs.
type ProcessInfo struct {
	// Image identifies the process image used to create the process.
	Image string
	// TODO(ericsnow) should be a list of strings?
	// Args is the extra args used to create the process.
	Args string
	// Desc is the description of the process.
	Desc string
	// Plugin identifies the process plugin used to manage the process.
	Plugin string

	// TODO(wwitzel3) determine specific details for storage and networking
	// Storage is the information used to identify storage resources
	// used when creating the process.
	Storage string
	// Networking is the information used to identify network resources
	// used when creating the process.
	Networking string

	// Details is the information about the process which the plugin provided.
	Details ProcessDetails
}

// TODO(wwitzel3) determine storageInfo based on spec/plugin needs
type storageInfo struct{}

// TODO(wwitzel3) determine networkInfo based on spec/plugin needs
type networkInfo struct{}

// PluginResource exposes an API for plugins to call
// to validate and allocate Juju resources.
type PluginResource interface {
	// Storage returns information needed by a plugin about the given
	// storage resource.
	// $ storage-info -> API.PluginResource.Storage
	Storage(storageID string) (storageInfo, error)
	// Networking returns information needed by a plugin about the given
	// networking resource.
	// $ network-info -> API.PluginResource.Networking
	Networking(networking string) (networkInfo, error)
}
