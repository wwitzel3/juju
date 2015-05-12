// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

// ProcManager defines the interface used to allow a process manager
// worker to perform actions.
type ProcManager interface {
	Launch() (Process, error)
	Destroy(UUID string) error

	List() ([]Process, error)
	Get(UUID string) (Process, error)
}

// Process contains the information about a specific process.
type Process struct {
	Image  string
	Args   string
	UUID   string
	Status string
	Type   string
}

// NewProcManager returns a worker which manages the
// life-cycle of a process for a given plugin type.
func NewProcManager(plugin string) ProcManager {
	return nil
}
