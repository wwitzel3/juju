// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package procmanager

import (
	"github.com/juju/juju/state"
)

// ProcManager defines the interface needed by the different areas of Juju.
type ProcManager interface {
	ProcManagerUnitAPI
	ProcManagerStateAPI
	ProcManagerWorker
}

// NewProcManager returns a
func NewProcManager(
	st *state.State,
	resources *common.Resources,
	authorizer common.Authorizer,
) (ProcManager, error) {
	return nil, nil
}

type procManager struct {
	UUID   string
	Status string
}
