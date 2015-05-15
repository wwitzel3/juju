package apiclient

// ProcessInfo holds the info about a process as passed through the API.
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

	// UniqueID is the identifier provided by the plugin when the
	// process was created.
	UniqueID string
	// Status is the most recent process status provided by the plugin.
	Status string
}

// ProcessID holds the info needed to identify a process through the API.
type ProcessID struct {
	// UUID is Juju's ID for the process.
	UUID string
}
