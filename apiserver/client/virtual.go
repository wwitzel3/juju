// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package client

import (
	"strings"

	"gopkg.in/juju/charm.v4"

	"github.com/juju/juju/apiserver/params"
	"github.com/juju/juju/state"
)

// VirtualServiceDeploy
func (c *Client) VirtualServiceDeploy(args params.VirtualServiceDeploy) error {
	vcharm := NewVirtualCharm(args)

	url := "virtual:virtual/" + args.ServiceName
	curl := charm.MustParseURL(url)
	stch, err := c.api.state.Charm(curl)
	if err != nil {
		stch, err = c.api.state.AddCharm(vcharm, curl, "", "global")
	}

	if err != nil {
		return err
	}

	env, err := c.api.state.Environment()
	if err != nil {
		return err
	}

	service, err := c.api.state.AddVirtualService(
		args.ServiceName,
		env.Owner().String(),
		stch,
	)
	if err != nil {
		return err
	}

	unit, err := service.AddVirtualUnit()
	if err != nil {
		return err
	}

	if err := setVirtualServiceSettings(c.api.state, args.ServiceName, args.Endpoints); err != nil {
		return err
	}

	if _, err := unit.SetAgentPresence(); err != nil {
		return err
	}
	if err := unit.SetAgentStatus(state.StatusActive, "virtual", nil); err != nil {
		return err
	}

	return nil
}

// makeVirtualRelations create relations map from virtual endpoints
func makeVirtualRelations(endpoints []params.VirtualEndpoint) map[string]charm.Relation {
	var relations = make(map[string]charm.Relation)
	for _, endpoint := range endpoints {
		relation := charm.Relation{
			Name:      endpoint.Relation,
			Role:      "provider",
			Interface: endpoint.Interface,
			Scope:     charm.ScopeGlobal,
		}
		relations[endpoint.Relation] = relation
	}
	return relations
}

func setVirtualServiceSettings(st *state.State, serviceName string, endpoints []params.VirtualEndpoint) error {
	for _, ep := range endpoints {
		key := strings.Join([]string{"s", serviceName, ep.Relation, ep.Interface}, "#")
		logger.Debugf("serviceKey: %q", key)
		state.WriteVirtualSettings(st, key, ep.Values)
	}
	return nil
}

type virtualCharm struct {
	meta *charm.Meta
}

var _ charm.Charm = (*virtualCharm)(nil)

// NewVirtualCharm returns a a virtual charm that is suitable for
// use with VirtualServiceDeploy.
func NewVirtualCharm(args params.VirtualServiceDeploy) virtualCharm {
	endpoints := makeVirtualRelations(args.Endpoints)
	meta := &charm.Meta{
		args.ServiceName,
		"",
		"",
		false,
		endpoints,
		nil,
		nil,
		0,
		0,
		nil,
		nil,
		"global",
		nil,
	}
	vcharm := virtualCharm{meta}
	return vcharm
}

func (vc virtualCharm) Meta() *charm.Meta {
	return vc.meta
}

func (virtualCharm) Actions() *charm.Actions {
	return nil
}

func (virtualCharm) Config() *charm.Config {
	return nil
}

func (virtualCharm) Metrics() *charm.Metrics {
	return nil
}

func (virtualCharm) Revision() int {
	return 0
}
