// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package environmentserver

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/state"
)

// distributeuUnit takes a unit and set of clean, possibly empty, instances
// and asks the InstanceDistributor policy (if any) which ones are suitable
// for assigning the unit to. If there is no InstanceDistributor, or the
// distribution group is empty, then all of the candidates will be returned.
func (d *deployer) DistributeUnit(u *state.Unit, candidates []instance.Id) ([]instance.Id, error) {
	if len(candidates) == 0 {
		return nil, nil
	}
	cfg, err := d.state.EnvironConfig()
	if err != nil {
		return nil, err
	}
	distributor, err := d.InstanceDistributor(cfg)
	if errors.IsNotImplemented(err) {
		return candidates, nil
	} else if err != nil {
		return nil, err
	}
	if distributor == nil {
		return nil, fmt.Errorf("policy returned nil instance distributor without an error")
	}
	distributionGroup, err := d.ServiceInstances(u.ServiceName())
	if err != nil {
		return nil, err
	}
	if len(distributionGroup) == 0 {
		return candidates, nil
	}
	return distributor.DistributeInstances(candidates, distributionGroup)
}

// ServiceInstances returns the instance IDs of provisioned
// machines that are assigned units of the specified service.
func (d *deployer) ServiceInstances(serviceName string) ([]instance.Id, error) {
	service, err := d.state.Service(serviceName)
	if err != nil {
		return nil, err
	}

	units, err := service.AllUnits()
	if err != nil {
		return nil, err
	}

	instanceIds := make([]instance.Id, 0, len(units))
	for _, unit := range units {
		machineId, err := unit.AssignedMachineId()
		if state.IsNotAssigned(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		machine, err := d.state.Machine(machineId)
		if err != nil {
			return nil, err
		}
		instanceId, err := machine.InstanceId()
		if err == nil {
			instanceIds = append(instanceIds, instanceId)
		} else if state.IsNotProvisionedError(err) {
			continue
		} else {
			return nil, err
		}
	}
	return instanceIds, nil
}