// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package status

import (
	"fmt"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/cmd/envcmd"
	"github.com/juju/juju/workload"
)

// ListAPI has the API methods needed by ListCommand.
type ListAPI interface {
	List(patterns ...string) ([]workload.Payload, error)
	Close() error
}

// ListCommand implements the list-payloads command.
type ListCommand struct {
	envcmd.EnvCommandBase
	out      cmd.Output
	patterns []string

	newAPIClient func(c *ListCommand) (ListAPI, error)
}

// NewListCommand returns a new command that lists charm payloads
// in the current environment.
func NewListCommand(newAPIClient func(c *ListCommand) (ListAPI, error)) *ListCommand {
	cmd := &ListCommand{
		newAPIClient: newAPIClient,
	}
	return cmd
}

var listDoc = `
This command will report on the runtime state of defined payloads.

Patterns can be one or more of:
- unit name
- machine id
- payload type
- payload class
- payload id
- payload tag
- payload status

When a pattern is specified, Juju will filter the status to only
those payloads that match their respective patterns.
`

func (c *ListCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "list-payloads",
		Args:    "[pattern ...]",
		Purpose: "display status information about currently running payloads",
		Doc:     listDoc,
	}
}

func (c *ListCommand) SetFlags(f *gnuflag.FlagSet) {
	defaultFormat := "tabular"
	c.out.AddFlags(f, defaultFormat, map[string]cmd.Formatter{
		"tabular": FormatTabular,
		"yaml":    cmd.FormatYaml,
		"json":    cmd.FormatJson,
	})
}

func (c *ListCommand) Init(args []string) error {
	c.patterns = args
	return nil
}

const connectionError = `Unable to connect to environment %q.
Please check your credentials or use 'juju bootstrap' to create a new environment.

Error details:
%v
`

func (c *ListCommand) Run(ctx *cmd.Context) error {
	apiclient, err := c.newAPIClient(c)
	if err != nil {
		return fmt.Errorf(connectionError, c.ConnectionName(), err)
	}
	defer apiclient.Close()

	infos, err := apiclient.List(c.patterns...)
	if err != nil {
		if infos == nil {
			// List call completely failed; there is nothing to report.
			return errors.Trace(err)
		}
		// Display any error, but continue to print info if some was returned.
		fmt.Fprintf(ctx.Stderr, "%v\n", err)
	} else if infos == nil {
		return errors.Errorf("unable to list the current payloads")
	}

	formatter := newListFormatter(infos, c.CompatVersion())
	formatted := formatter.format()
	return c.out.Write(ctx, formatted)
}
