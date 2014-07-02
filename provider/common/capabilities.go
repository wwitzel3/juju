// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"github.com/juju/juju/environs/instances"
	"github.com/juju/juju/state/api/params"
)

type EnvironCapability interface {
	AvailabilityZones() ([]AvailabilityZone, error)
	InstanceTypes() ([]instances.InstanceType, error)
	SupportedArchitectures() ([]string, error)
}

type EnvironCapabilities struct {
	AvailabilityZones      []AvailabilityZone
	InstanceTypes          []instances.InstanceType
	SupportedArchitectures []string
}

func NewEnvironCapabilities(env EnvironCapability) (*EnvironCapabilities, error) {
	availabilityZones, err := env.AvailabilityZones()
	if err != nil {
		return nil, err
	}

	instanceTypes, err := env.InstanceTypes()
	if err != nil {
		return nil, err
	}

	supportedArchitectures, err := env.SupportedArchitectures()
	if err != nil {
		return nil, err
	}

	return &EnvironCapabilities{
		AvailabilityZones:      availabilityZones,
		InstanceTypes:          instanceTypes,
		SupportedArchitectures: supportedArchitectures,
	}, nil
}

func (ec *EnvironCapabilities) Result() params.EnvironmentCapabilitiesResult {
	return params.EnvironmentCapabilitiesResult{
		AvailabilityZones:      ec.encodeAvailabilityZones(),
		InstanceTypes:          ec.encodeInstanceTypes(),
		SupportedArchitectures: ec.SupportedArchitectures,
	}
}

func (ec *EnvironCapabilities) encodeAvailabilityZones() []map[string]interface{} {
	availabilityZones := make([]map[string]interface{}, len(ec.AvailabilityZones), len(ec.AvailabilityZones))
	for i, availabilityZone := range ec.AvailabilityZones {
		availabilityZones[i] = map[string]interface{}{
			"name":      availabilityZone.Name(),
			"available": availabilityZone.Available(),
		}
	}
	return availabilityZones
}

func (ec *EnvironCapabilities) encodeInstanceTypes() []map[string]interface{} {
	instanceTypes := make([]map[string]interface{}, len(ec.InstanceTypes))
	for i, instanceType := range ec.InstanceTypes {
		instanceTypes[i] = map[string]interface{}{
			"id":       instanceType.Id,
			"name":     instanceType.Name,
			"cpucores": instanceType.CpuCores,
			"memory":   instanceType.Mem,
			"disk":     instanceType.RootDisk,
		}
	}
	return instanceTypes
}
