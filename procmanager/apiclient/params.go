package apiclient

type ProcessInfo struct {
	Image  string
	Args   string
	Desc   string
	Plugin string

	// TODO(wwitzel3) determine specific details for storage and networking
	Storage    string
	Networking string

	// ProcessDetails
	UniqueID string
	Status   string
}

type ProcessID struct {
	UUID string
}
