// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"github.com/juju/cmd"

	"github.com/juju/juju/worker/uniter/runner/jujuc"
)

func init() {
	jujuc.RegisterCommand("destroy", NewDestroyCommand)
}

// DestroyCommand implements the destroy command.
type DestroyCommand struct {
	cmd.CommandBase
	ctx        jujuc.Context
	RelationId int
	Key        string
	UnitName   string
	out        cmd.Output
}

// NewDestroyCommand returns a new DestroyCommand.
func NewDestroyCommand(ctx jujuc.Context) cmd.Command {
	return &DestroyCommand{ctx: ctx}
}

// Info implements cmd.Command.Info.
func (c *DestroyCommand) Info() *cmd.Info {
	args := "<uuid>"
	// TODO(ericsnow) finish
	doc := `
destroy prints the value of a unit's relation setting, specified by key.
If no key is given, or if the key is "-", all keys and values will be printed.
`
	return &cmd.Info{
		Name: "destroy",
		Args: args,
		// TODO(ericsnow) finish
		Purpose: "get relation settings",
		Doc:     doc,
	}
}

// Init implements cmd.Command.Init.
func (c *DestroyCommand) Init(args []string) error {
	// TODO(ericsnow) finish
	return cmd.CheckEmpty(args)
}

// Run implements cmd.Command.Run.
func (c *DestroyCommand) Run(ctx *cmd.Context) error {
	// TODO(ericsnow) finish
	// $ destroy -> API.CharmHookEnv.Destroy
	// valid arg parse
	// verify UUID
	// exec plugin
	// handle/surface errors
	// unregister UUID with state
	return c.out.Write(ctx, nil)
}
