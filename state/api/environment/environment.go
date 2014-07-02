// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environment

import (
	"github.com/juju/juju/state/api/base"
	"github.com/juju/juju/state/api/common"
	"github.com/juju/juju/state/api/params"
)

const apiName = "Environment"

// Facade provides access to a machine environment worker's view of the world.
type Facade struct {
	*common.EnvironWatcher
	caller base.Caller
}

// NewFacade returns a new api client facade instance.
func NewFacade(caller base.Caller) *Facade {
	return &Facade{
		EnvironWatcher: common.NewEnvironWatcher(apiName, caller),
		caller:         caller,
	}
}

func (f *Facade) GetCapabilities() (params.EnvironmentCapabilitiesResult, error) {
	var results params.EnvironmentCapabilitiesResult

	err := f.caller.Call(apiName, "", "GetCapabilities", nil, &results)
	if err != nil {
		return results, err
	}

	return results, nil
}
