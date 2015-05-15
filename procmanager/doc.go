// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

/*
procmanager exposes the ability for charm authors to hand off the
management of creating and destroying processes to juju. By handing the
creation and destruction of these processes you enable juju to surface
these running processes to viewers of a unit's status, giving the viewer
a more accurate description of the environment.

TODO(ericsnow) more info here?

The following PlantUML diagram describes the interactions between
different components when a charm launches a new process.

  TODO(wwitzel3) INSERT PLANTUML HERE
*/
package procmanager
