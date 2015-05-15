// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package all

import (
	// Ensure we register the procmanager API facade.
	_ "github.com/juju/juju/procmanager/apiserver"
	// Ensure we register procmanager commands.
	_ "github.com/juju/juju/procmanager/commands"
)
