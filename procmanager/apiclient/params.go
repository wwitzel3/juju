package apiclient

type AddArgs struct {
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

type RemoveArgs struct {
}

type InfoArgs struct {
}

type InfoResult struct {
}
