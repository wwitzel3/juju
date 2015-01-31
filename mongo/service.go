// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package mongo

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"

	"github.com/juju/errors"
	"github.com/juju/utils"

	"github.com/juju/juju/service"
	"github.com/juju/juju/service/common"
)

const (
	maxFiles = 65000
	maxProcs = 20000

	serviceName = "juju-db"
)

// TODO(ericsnow) Move JujuMongodPath over to files.go?

// These are vars for the sake of patching during tests.
var (
	// JujuMongodPath holds the default path to the juju-specific mongod.
	JujuMongodPath = "/usr/lib/juju/bin/mongod"
)

// These constants relate to MongoDB Numa support.
const (
	// This is the name of the variable to use in ExtraScript
	// fragment to substitute into init script template.
	multinodeVarName = "MULTI_NODE"
	// This value will be used to wrap desired mongo cmd in numactl if wanted/needed
	numaCtlWrap = "$%v"
	// Extra shell script fragment for init script template.
	// This determines if we are dealing with multi-node environment
	detectMultiNodeScript = `%v=""
if [ $(find /sys/devices/system/node/ -maxdepth 1 -mindepth 1 -type d -name node\* | wc -l ) -gt 1 ]
then
    %v=" numactl --interleave=all "
    # Ensure sysctl turns off zone_reclaim_mode if not already set
    (grep -q vm.zone_reclaim_mode /etc/sysctl.conf || echo vm.zone_reclaim_mode = 0 >> /etc/sysctl.conf) && sysctl -p
fi
`
)

// Path returns the executable path to be used to run mongod on this
// machine. If the juju-bundled version of mongo exists, it will return that
// path, otherwise it will return the command to run mongod from the path.
func Path() (string, error) {
	if _, err := os.Stat(JujuMongodPath); err == nil {
		return JujuMongodPath, nil
	}

	path, err := exec.LookPath("mongod")
	if err != nil {
		logger.Infof("could not find %v or mongod in $PATH", JujuMongodPath)
		return "", err
	}
	return path, nil
}

// ServiceName returns the name of the init service config for mongo using
// the given namespace.
func ServiceName(namespace string) string {
	if namespace != "" {
		return fmt.Sprintf("%s-%s", serviceName, namespace)
	}
	return serviceName
}

// ServiceSpec holds all the information necessary to interact with a
// juju-managed mongodb service.
type ServiceSpec struct {
	Executable  string
	DBDir       string
	DataDir     string
	Port        int
	OplogSizeMB int
	WantNumaCtl bool
}

// ApplyDefaults sets any unset fields to their correct defaults.
func (ss *ServiceSpec) ApplyDefaults() error {
	if ss.Executable == "" {
		mongoPath, err := Path()
		if err != nil {
			return errors.Trace(err)
		}
		ss.Executable = mongoPath
		logVersion(mongoPath)
	}

	return nil
}

func (ss ServiceSpec) command() string {
	// ss.Executable must be set (call ss.ApplyDefaults if necessary).
	return ss.Executable +
		" --auth" +
		" --dbpath=" + utils.ShQuote(ss.DBDir) +
		" --sslOnNormalPorts" +
		" --sslPEMKeyFile " + utils.ShQuote(sslKeyPath(ss.DataDir)) +
		" --sslPEMKeyPassword ignored" +
		" --port " + fmt.Sprint(ss.Port) +
		" --noprealloc" +
		" --syslog" +
		" --smallfiles" +
		" --journal" +
		" --keyFile " + utils.ShQuote(sharedSecretPath(ss.DataDir)) +
		" --replSet " + ReplicaSetName +
		" --ipv6 " +
		" --oplogSize " + strconv.Itoa(ss.OplogSizeMB)
}

// Conf builds a new service.Conf from the spec and returns it.
func (ss ServiceSpec) Conf() common.Conf {
	mongoCmd := ss.command()

	extraScript := ""
	if ss.WantNumaCtl {
		extraScript = fmt.Sprintf(detectMultiNodeScript, multinodeVarName, multinodeVarName)
		mongoCmd = fmt.Sprintf(numaCtlWrap, multinodeVarName) + mongoCmd
	}

	conf := common.Conf{
		Desc: "juju state database",
		Cmd:  mongoCmd,
		Limit: map[string]string{
			"nofile": fmt.Sprintf("%d %d", maxFiles, maxFiles),
			"nproc":  fmt.Sprintf("%d %d", maxProcs, maxProcs),
		},
		ExtraScript: extraScript,
	}
	return conf
}

var newService = func(name, dataDir string, conf common.Conf) (*service.Service, error) {
	return service.NewService(name, dataDir, conf)
}

// NewService builds a new service based on the spec and returns it.
func (ss ServiceSpec) NewService(namespace string) (*Service, error) {
	name := ServiceName(namespace)
	svc, err := newService(name, ss.DataDir, ss.Conf())
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Service{*svc}, nil
}

// Service represents the a juju-managed mongodb service.
type Service struct {
	service.Service
}

func (svc Service) startIfInstalled() (bool, error) {
	err := svc.Manage()
	if err != nil && !errors.IsAlreadyExists(err) {
		return false, errors.Trace(err)
	}

	enabled, err := svc.IsEnabled()
	if err != nil {
		return false, errors.Trace(err)
	}
	if !enabled {
		return false, nil
	}

	matched, err := svc.Check()
	if err != nil {
		return false, errors.Trace(err)
	}
	if !matched {
		return false, nil
	}

	logger.Debugf("mongo exists as expected")
	running, err := svc.IsRunning()
	if err != nil {
		return false, errors.Trace(err)
	}

	if !running {
		if err := svc.Start(); err != nil {
			return false, errors.Trace(err)
		}
	}

	return true, nil
}

// noauthCommand returns an os/exec.Cmd that may be executed to
// run mongod without security.
func noauthCommand(dataDir string, port int) (*exec.Cmd, error) {
	sslKeyFile := path.Join(dataDir, "server.pem")
	dbDir := filepath.Join(dataDir, "db")
	mongoPath, err := Path()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(mongoPath,
		"--noauth",
		"--dbpath", dbDir,
		"--sslOnNormalPorts",
		"--sslPEMKeyFile", sslKeyFile,
		"--sslPEMKeyPassword", "ignored",
		"--bind_ip", "127.0.0.1",
		"--port", fmt.Sprint(port),
		"--noprealloc",
		"--syslog",
		"--smallfiles",
		"--journal",
	)
	return cmd, nil
}
