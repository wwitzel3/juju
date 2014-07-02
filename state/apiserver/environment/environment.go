// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environment

import (
	"fmt"
	"github.com/juju/juju/environs"
	providerCommon "github.com/juju/juju/provider/common"
	"github.com/juju/juju/state"
	"github.com/juju/juju/state/api/params"
	"github.com/juju/juju/state/apiserver/common"
)

func init() {
	common.RegisterStandardFacade("Environment", 0, NewEnvironmentAPI)
}

// EnvironmentAPI implements the API used by the machine environment worker.
type EnvironmentAPI struct {
	*common.EnvironWatcher
	st *state.State
}

// NewEnvironmentAPI creates a new instance of the Environment API.
func NewEnvironmentAPI(st *state.State, resources *common.Resources, authorizer common.Authorizer) (*EnvironmentAPI, error) {
	// Can always watch for environ changes.
	getCanWatch := common.AuthAlways(true)
	// Does not get the secrets.
	getCanReadSecrets := common.AuthAlways(false)
	return &EnvironmentAPI{
		EnvironWatcher: common.NewEnvironWatcher(st, resources, getCanWatch, getCanReadSecrets),
		st:             st,
	}, nil
}

func (api *EnvironmentAPI) GetCapabilities() (params.EnvironmentCapabilitiesResult, error) {
	emptyEnvironmentCapabilities := params.EnvironmentCapabilitiesResult{}

	environConfig, err := api.st.EnvironConfig()
	if err != nil {
		return emptyEnvironmentCapabilities, err
	}

	environment, err := environs.New(environConfig)
	if err != nil {
		return emptyEnvironmentCapabilities, err
	}

	capabilityer, ok := environment.(providerCommon.EnvironCapability)
	if !ok {
		return emptyEnvironmentCapabilities, fmt.Errorf("environment does not implement the correct interface: %v", ok)
	}

	environCapabilities, err := providerCommon.NewEnvironCapabilities(capabilityer)
	if err != nil {
		return emptyEnvironmentCapabilities, err
	}

	return environCapabilities.Result(), nil
}
