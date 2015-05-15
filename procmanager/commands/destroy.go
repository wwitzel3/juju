// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package commands

import (
	"fmt"

	"github.com/juju/cmd"

	"github.com/juju/juju/apiserver/params"
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

func NewDestroyCommand(ctx jujuc.Context) cmd.Command {
	return &DestroyCommand{ctx: ctx}
}

func (c *DestroyCommand) Info() *cmd.Info {
	args := "<uuid>"
	doc := `
destroy prints the value of a unit's relation setting, specified by key.
If no key is given, or if the key is "-", all keys and values will be printed.
`
	return &cmd.Info{
		Name:    "destroy",
		Args:    args,
		Purpose: "get relation settings",
		Doc:     doc,
	}
}

func (c *DestroyCommand) Init(args []string) error {
	if c.RelationId == -1 {
		return fmt.Errorf("no relation id specified")
	}
	c.Key = ""
	if len(args) > 0 {
		if c.Key = args[0]; c.Key == "-" {
			c.Key = ""
		}
		args = args[1:]
	}
	if name, found := c.ctx.RemoteUnitName(); found {
		c.UnitName = name
	}
	if len(args) > 0 {
		c.UnitName = args[0]
		args = args[1:]
	}
	if c.UnitName == "" {
		return fmt.Errorf("no unit id specified")
	}
	return cmd.CheckEmpty(args)
}

func (c *DestroyCommand) Run(ctx *cmd.Context) error {
	r, found := c.ctx.Relation(c.RelationId)
	if !found {
		return fmt.Errorf("unknown relation id")
	}
	var settings params.Settings
	if c.UnitName == c.ctx.UnitName() {
		node, err := r.Settings()
		if err != nil {
			return err
		}
		settings = node.Map()
	} else {
		var err error
		settings, err = r.ReadSettings(c.UnitName)
		if err != nil {
			return err
		}
	}
	if c.Key == "" {
		return c.out.Write(ctx, settings)
	}
	if value, ok := settings[c.Key]; ok {
		return c.out.Write(ctx, value)
	}
	return c.out.Write(ctx, nil)
}
