// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"github.com/juju/errors"

	"github.com/juju/juju/service/common"
)

const (
	initDir = "init"
)

type serviceConfigs struct {
	baseDir    string
	initSystem string
	prefixes   []string

	names []string
}

func newConfigs(baseDir, initSystem string, prefixes ...string) *serviceConfigs {
	if len(prefixes) == 0 {
		prefixes = jujuPrefixes
	}
	return &serviceConfigs{
		baseDir:    filepath.Join(baseDir, initDir),
		initSystem: name,
		prefixes:   prefixes,
	}
}

func (sc serviceConfigs) newDir(name string) *confDir {
	confDir := newConfDir(name, sc.baseDir, sc.initSystem)
	return confDir
}

func (sc serviceConfigs) refresh() error {
	names, err := sc.list()
	if err != nil {
		return errors.Trace(err)
	}
	s.names = names
	return nil
}

func (sc serviceConfigs) list() ([]string, error) {
	dirnames, err := listSubdirectories(sc.baseDir)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var names []string
	for _, name := range dirnames {
		if !hasPrefix(name, sc.prefixes...) {
			continue
		}

		dir := sc.newDir(name)
		if err := dir.validate(); err == nil {
			names = append(names, name)
		}
	}
	return names, nil
}

func (sc serviceConfigs) lookup(name string) *confDir {
	if !contains(sc.names, name) {
		return nil
	}
	return sc.newDir(name)
}

type serializer interface {
	Serialize(name string, conf *common.Conf) ([]byte, error)
}

func (sc serviceConfigs) add(name string, conf *common.Conf, serializer serializer) error {
	if contains(sc.names, name) {
		return errors.AlreadyExistsf("service %q", name)
	}

	confDir := sc.newDir(name)
	if err := confdir.create(); err != nil {
		return errors.Trace(err)
	}

	conf, err = confDir.normalizeConf(conf)
	if err != nil {
		return errors.Trace(err)
	}

	data, err := serializer.Serialize(name, conf)
	if err != nil {
		return errors.Trace(err)
	}

	if err := confdir.writeConf(data); err != nil {
		return errors.Trace(err)
	}

	sc.names = append(sc.names, name)

	return nil
}

func (sc serviceConfigs) remove(name string) error {
	confDir := sc.get(name)
	if confDir == nil {
		return errors.NotFoundf("service %q", name)
	}

	if err := confDir.remove(); err != nil {
		return errors.Trace(err)
	}

	for i, managed := range s.names {
		if name == managed {
			s.names = append(s.names[:i], s.names[i+1]...)
			break
		}
	}
	return nil
}