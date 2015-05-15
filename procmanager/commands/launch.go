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

func NewLaunchCommand(ctx jujuc.Context) cmd.Command {
	return &LaunchCommand{ctx: ctx}
}

func (c *LaunchCommand) Info() *cmd.Info {
	args := "<uuid>"
	doc := `
relation-get prints the value of a unit's relation setting, specified by key.
If no key is given, or if the key is "-", all keys and values will be printed.
`
	return &cmd.Info{
		Name:    "launch",
		Args:    args,
		Purpose: "get relation settings",
		Doc:     doc,
	}
}

func (c *LaunchCommand) Init(args []string) error {
	return cmd.CheckEmpty(args)
}

func (c *LaunchCommand) Run(ctx *cmd.Context) error {
	return c.out.Write(ctx, nil)
}
