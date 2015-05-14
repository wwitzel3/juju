/* procmanager exposes the ability for charm authors to hand off the management of processes
 * created by the charm to juju. By handing the creation and destrouction of these processes
 * you enable juju to surface these running processes to viewers of a units status, giving
 * the viewer a more accurate description of the environment.
 *
 * procmanager exposes a single interface `ProcManager` that juju can use to perform all of
 * the process management tasks, below is a PlantUML diagram that shows the interaction of
 * the procmanager module and juju.
 *
 *  INSERT PLANTUML HERE
 */
package procmanager
