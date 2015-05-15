// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

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
