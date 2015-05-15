// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package worker

import (
	"github.com/juju/juju/cmd/jujud/agent"
	"github.com/juju/juju/worker"
)

func init() {
	if err := agent.RegisterWorker("proc-manager", newWorker); err != nil {
		panic(err)
	}
}

// ProcManagerUnitAPI
type ProcManagerWorker interface {
	Handle()
}

func newWorker() (worker.Worker, error) {
	worker := worker.Worker(nil)

	return worker, nil
}
