// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/juju/cmd"
	"github.com/juju/errors"
	"github.com/juju/names"
	"launchpad.net/gnuflag"

	"github.com/juju/juju/api/runcmd"
	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/cmd/envcmd"
)

// RunCommand is responsible for running arbitrary commands on remote machines.
type RunCommand struct {
	envcmd.EnvCommandBase
	out        cmd.Output
	all        bool
	timeout    time.Duration
	machines   []string
	services   []string
	units      []string
	relation   string
	remoteUnit string
	commands   string
}

const runDoc = `
Run the commands on the specified targets.

Targets are specified using either machine ids, service names or unit
names.  At least one target specifier is needed.

Multiple values can be set for --machine, --service, and --unit by using
comma separated values.

If the target is a machine, the command is run as the "ubuntu" user on
the remote machine.

If the target is a service, the command is run on all units for that
service. For example, if there was a service "mysql" and that service
had two units, "mysql/0" and "mysql/1", then
  --service mysql
is equivalent to
  --unit mysql/0,mysql/1

Commands run for services or units are executed in a 'hook context' for
the unit.

--relation allows you to ensure the command is executed on the specified
service or unit targets with a specific relation context.

--remote-unit is used with --relation to specify a remote-unit in cases where
more than one exists. If only one remote-unit exists there is no need to specify this.

--all is provided as a simple way to run the command on all the machines
in the environment.  If you specify --all you cannot provide additional
targets or a relation.

`

func (c *RunCommand) Info() *cmd.Info {
	return &cmd.Info{
		Name:    "run",
		Args:    "<commands>",
		Purpose: "run the commands on the remote targets specified",
		Doc:     runDoc,
	}
}

func (c *RunCommand) SetFlags(f *gnuflag.FlagSet) {
	c.out.AddFlags(f, "smart", cmd.DefaultFormatters)
	f.BoolVar(&c.all, "all", false, "run the commands on all the machines")
	f.DurationVar(&c.timeout, "timeout", 5*time.Minute, "how long to wait before the remote command is considered to have failed")
	f.Var(cmd.NewStringsValue(nil, &c.machines), "machine", "one or more machine ids")
	f.Var(cmd.NewStringsValue(nil, &c.services), "service", "one or more service names")
	f.Var(cmd.NewStringsValue(nil, &c.units), "unit", "one or more unit ids")
	f.StringVar(&c.relation, "relation", "", "relation context to run the command under")
	f.StringVar(&c.remoteUnit, "remote-unit", "", "run the command for a specific remote unit for a given relation")
}

func (c *RunCommand) Init(args []string) error {
	if len(args) == 0 {
		return errors.Errorf("no commands specified")
	}
	c.commands, args = args[0], args[1:]

	if c.all {
		if len(c.machines) != 0 {
			return errors.Errorf("You cannot specify --all and individual machines")
		}
		if len(c.services) != 0 {
			return errors.Errorf("You cannot specify --all and individual services")
		}
		if len(c.units) != 0 {
			return errors.Errorf("You cannot specify --all and individual units")
		}
		if len(c.relation) != 0 {
			return errors.Errorf("You cannot specify --all and a relation")
		}
		if len(c.remoteUnit) != 0 {
			return errors.Errorf("You cannot specify --all and a remote-unit")
		}

	} else {
		if len(c.machines) == 0 && len(c.services) == 0 && len(c.units) == 0 {
			return errors.Errorf("You must specify a target, either through --all, --machine, --service or --unit")
		}
		if len(c.relation) == 0 && len(c.remoteUnit) != 0 {
			return errors.Errorf("You must specify a relation through --relation")
		}
	}

	if len(c.machines) != 0 {
		if len(c.relation) != 0 {
			return errors.Errorf("You cannot specify --machine and a relations")
		}
		if len(c.remoteUnit) != 0 {
			return errors.Errorf("You cannot specify --machine and a remote-unit")
		}
	}

	var nameErrors []string
	for _, machineId := range c.machines {
		if !names.IsValidMachine(machineId) {
			nameErrors = append(nameErrors, fmt.Sprintf("  %q is not a valid machine id", machineId))
		}
	}
	for _, service := range c.services {
		if !names.IsValidService(service) {
			nameErrors = append(nameErrors, fmt.Sprintf("  %q is not a valid service name", service))
		}
	}
	for _, unit := range c.units {
		if !names.IsValidUnit(unit) {
			nameErrors = append(nameErrors, fmt.Sprintf("  %q is not a valid unit name", unit))
		}
	}

	if len(c.remoteUnit) > 0 && !names.IsValidUnit(c.remoteUnit) {
		nameErrors = append(nameErrors, fmt.Sprintf("  %q is not a valid remote-unit name", c.remoteUnit))
	}

	if len(nameErrors) > 0 {
		return errors.Errorf("The following run targets are not valid:\n%s",
			strings.Join(nameErrors, "\n"))
	}

	return cmd.CheckEmpty(args)
}

func encodeBytes(input []byte) (value string, encoding string) {
	if utf8.Valid(input) {
		value = string(input)
		encoding = "utf8"
	} else {
		value = base64.StdEncoding.EncodeToString(input)
		encoding = "base64"
	}
	return value, encoding
}

func storeOutput(values map[string]interface{}, key string, input []byte) {
	value, encoding := encodeBytes(input)
	values[key] = value
	if encoding != "utf8" {
		values[key+".encoding"] = encoding
	}
}

// ConvertRunResults takes the results from the api and creates a map
// suitable for format converstion to YAML or JSON.
func ConvertRunResults(runResults []params.RunResult) interface{} {
	var results = make([]interface{}, len(runResults))

	for i, result := range runResults {
		// We always want to have a string for stdout, but only show stderr,
		// code and error if they are there.
		values := make(map[string]interface{})
		values["MachineId"] = result.MachineId
		if result.UnitId != "" {
			values["UnitId"] = result.UnitId

		}
		storeOutput(values, "Stdout", result.Stdout)
		if len(result.Stderr) > 0 {
			storeOutput(values, "Stderr", result.Stderr)
		}
		if result.Code != 0 {
			values["ReturnCode"] = result.Code
		}
		if result.Error != "" {
			values["Error"] = result.Error
		}
		results[i] = values
	}

	return results
}

func (c *RunCommand) Run(ctx *cmd.Context) error {
	root, err := c.NewAPIRoot()
	if err != nil {
		return errors.Annotate(err, "cannot get API connection")
	}

	runClient := runcmd.NewClient(root, root.EnvironTag())
	defer runClient.Close()

	var runResults []params.RunResult
	params := params.RunParams{
		Commands:   c.commands,
		Timeout:    c.timeout,
		Machines:   c.machines,
		Services:   c.services,
		Units:      c.units,
		Relation:   c.relation,
		RemoteUnit: c.remoteUnit,
	}

	if c.all {
		runResults, err = runClient.RunOnAllMachines(params.Commands, params.Timeout)
	} else {
		runResults, err = runClient.Run(params)
	}

	if err != nil {
		oldClient, err := getRunAPIClient(c)
		if err != nil {
			return errors.Annotate(err, "unable to get a suitable client")
		}

		if c.all {
			runResults, err = oldClient.RunOnAllMachines(params.Commands, params.Timeout)
		} else {
			if len(params.Relation) > 0 || len(params.RemoteUnit) > 0 {
				return errors.Errorf("option(s) --relation, --remote-unit are not supported by this server")
			}
			runResults, err = oldClient.Run(params)
		}

		if err != nil {
			return err
		}
	}

	// If we are just dealing with one result, AND we are using the smart
	// format, then pretend we were running it locally.
	if len(runResults) == 1 && c.out.Name() == "smart" {
		result := runResults[0]
		ctx.Stdout.Write(result.Stdout)
		ctx.Stderr.Write(result.Stderr)
		if result.Error != "" {
			// Convert the error string back into an error object.
			return errors.Errorf("%s", result.Error)
		}
		if result.Code != 0 {
			return cmd.NewRcPassthroughError(result.Code)
		}
		return nil
	}

	c.out.Write(ctx, ConvertRunResults(runResults))
	return nil
}

// In order to be able to easily mock out the API side for testing,
// the API client is got using a function.

type RunClient interface {
	Close() error
	RunOnAllMachines(commands string, timeout time.Duration) ([]params.RunResult, error)
	Run(run params.RunParams) ([]params.RunResult, error)
}

// Here we need the signature to be correct for the interface.
var getRunAPIClient = func(c *RunCommand) (RunClient, error) {
	return c.NewAPIClient()
}
