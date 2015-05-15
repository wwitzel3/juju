// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package worker

import (
	"github.com/juju/juju/api"
	"github.com/juju/juju/cmd/jujud/agent"
)

func init() {
	err := agent.RegisterSimpleWorker("proc-watcher",
		func(st *api.State) (func(<-chan struct{}) error, error) {
			return newWatcherLoop(st)
		},
	)
	if err != nil {
		panic(err)
	}
}

func newWatcherLoop(st *api.State) (func(<-chan struct{}) error, error) {
	// TODO(ericsnow) finish
	return nil, nil
}
