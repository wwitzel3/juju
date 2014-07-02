package common

import (
	"github.com/juju/juju/environs/instances"
)

type InstanceTyper interface {
	InstanceTypes() ([]instances.InstanceType, error)
}

func InstanceTypeNames(env InstanceTyper) ([]string, error) {
	instanceTypes, err := env.InstanceTypes()
	if err != nil {
		return nil, err
	}
	instanceTypeNames := make([]string, len(instanceTypes))
	for i, instanceType := range instanceTypes {
		instanceTypeNames[i] = instanceType.Name
	}
	return instanceTypeNames, nil
}
