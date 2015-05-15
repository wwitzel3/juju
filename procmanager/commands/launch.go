// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"github.com/juju/cmd"

	"github.com/juju/juju/worker/uniter/runner/jujuc"
)

func init() {
	jujuc.RegisterCommand("launch", NewLaunchCommand)
}

// LaunchCommand implements the launch command.
type LaunchCommand struct {
	cmd.CommandBase
	ctx    jujuc.Context
	Plugin string
	out    cmd.Output
}

// NewLaunchCommand returns a new LaunchCommand.
func NewLaunchCommand(ctx jujuc.Context) cmd.Command {
	return &LaunchCommand{ctx: ctx}
}

// Info implements cmd.Command.Info.
func (c *LaunchCommand) Info() *cmd.Info {
	args := "<uuid>"
	// TODO(ericsnow) finish
	doc := `
relation-get prints the value of a unit's relation setting, specified by key.
If no key is given, or if the key is "-", all keys and values will be printed.
`
	return &cmd.Info{
		Name: "launch",
		Args: args,
		// TODO(ericsnow) finish
		Purpose: "get relation settings",
		Doc:     doc,
	}
}

// Init implements cmd.Command.Init.
func (c *LaunchCommand) Init(args []string) error {
	// TODO(ericsnow) finish
	return cmd.CheckEmpty(args)
}

// Run implements cmd.Command.Run.
func (c *LaunchCommand) Run(ctx *cmd.Context) error {
	// TODO(ericsnow) finish
	// $ launch -> API.CharmHookEnv.Launch
	// valid arg parse
	// exec plugin with validated arguments
	//     -> plugin will use PluginResources to verify storage and networking
	// handle any plugin errors
	// convert unique identifier in to UUID for process
	// register UUID and process information with state
	return c.out.Write(ctx, nil)
}
