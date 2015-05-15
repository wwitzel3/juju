// Copyright 2012-2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

// all triggers the import of all juju modules to ensure that any
// init methods for a module are executed.
package all

import (
	_ "github.com/juju/juju/procmanager"
)
