// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package initsystems

import (
	"github.com/juju/errors"
)

type enabledChecker interface {
	Info(name string) (*ServiceInfo, error)
}

// EnsureStatus may be used by InitSystem implementations to ensure
// that the named service has been enabled and the current status matches
// the provided status. Note: An empty string is treated as StatusEnabled.
// This function is important for operations where the service must first
// be enabled.
func EnsureStatus(is enabledChecker, name string, status string) error {
	info, err := is.Info(name)
	if status == StatusDisabled {
		if errors.IsNotFound(err) {
			return nil
		}
		return errors.AlreadyExistsf("service %q", name)
	}

	if err != nil {
		return errors.Trace(err)
	}

	if status == StatusEnabled {
		return nil
	}

	if info.Status == status {
		return nil
	}

	switch status {
	case StatusRunning:
		err = errors.NotFoundf("service %q", name)
	case StatusStopped:
		err = errors.AlreadyExistsf("service %q", name)
	default:
		err = errors.NotFoundf("service %q", name)
	}

	return err
}

// FilterNames filters out any name in names that isn't in include.
func FilterNames(names, include []string) []string {
	if len(include) == 0 {
		return names
	}

	var filtered []string
	for _, name := range names {
		for _, included := range include {
			if name == included {
				filtered = append(filtered, name)
				break
			}
		}
	}
	return filtered
}
