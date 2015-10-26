// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package lxd

import (
	"github.com/juju/errors"
	gitjujutesting "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/arch"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/cloudconfig/instancecfg"
	"github.com/juju/juju/cloudconfig/providerinit"
	"github.com/juju/juju/constraints"
	"github.com/juju/juju/container/lxd/lxdclient"
	"github.com/juju/juju/environs"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/network"
	"github.com/juju/juju/testing"
	"github.com/juju/juju/tools"
	"github.com/juju/juju/version"
)

// These values are stub LXD client credentials for use in tests.
const (
	PublicKey = `-----BEGIN CERTIFICATE-----
...
...
...
...
...
...
...
...
...
...
...
...
...
...
-----END CERTIFICATE-----
`
	PrivateKey = `-----BEGIN PRIVATE KEY-----
...
...
...
...
...
...
...
...
...
...
...
...
...
...
-----END PRIVATE KEY-----
`
)

// These are stub config values for use in tests.
var (
	ConfigAttrs = testing.FakeConfig().Merge(testing.Attrs{
		"type":        "lxd",
		"namespace":   "",
		"remote":      "",
		"client-cert": "",
		"client-key":  "",
		"uuid":        "2d02eeac-9dbb-11e4-89d3-123b93f75cba",
	})
)

// We test these here since they are not exported.
var (
	_ environs.Environ  = (*environ)(nil)
	_ instance.Instance = (*environInstance)(nil)
)

type BaseSuiteUnpatched struct {
	gitjujutesting.IsolationSuite

	Config    *config.Config
	EnvConfig *environConfig
	Env       *environ
	Prefix    string

	Addresses     []network.Address
	Instance      *environInstance
	RawInstance   *lxdclient.Instance
	InstName      string
	Hardware      *lxdclient.InstanceHardware
	HWC           *instance.HardwareCharacteristics
	Metadata      map[string]string
	StartInstArgs environs.StartInstanceParams
	//InstanceType  instances.InstanceType

	Ports []network.PortRange
}

func (s *BaseSuiteUnpatched) SetUpTest(c *gc.C) {
	s.IsolationSuite.SetUpTest(c)

	s.initEnv(c)
	s.initInst(c)
	s.initNet(c)
}

func (s *BaseSuiteUnpatched) initEnv(c *gc.C) {
	s.Env = &environ{
		name: "lxd",
	}
	cfg := s.NewConfig(c, nil)
	s.setConfig(c, cfg)
}

func (s *BaseSuiteUnpatched) initInst(c *gc.C) {
	tools := []*tools.Tools{{
		Version: version.Binary{Arch: arch.AMD64, Series: "trusty"},
		URL:     "https://example.org",
	}}

	cons := constraints.Value{
	// nothing
	}

	instanceConfig, err := instancecfg.NewBootstrapInstanceConfig(cons, "trusty")
	c.Assert(err, jc.ErrorIsNil)

	instanceConfig.Tools = tools[0]
	instanceConfig.AuthorizedKeys = s.Config.AuthorizedKeys()

	userData, err := providerinit.ComposeUserData(instanceConfig, nil, lxdRenderer{})
	c.Assert(err, jc.ErrorIsNil)

	s.Hardware = &lxdclient.InstanceHardware{
		Architecture: arch.AMD64,
		NumCores:     1,
		MemoryMB:     3750,
		//RootDiskMB:   50000,
	}
	var archName string = arch.AMD64
	var numCores uint64 = 1
	var memoryMB uint64 = 3750
	//var rootDiskMB uint64 = 50000
	s.HWC = &instance.HardwareCharacteristics{
		Arch:     &archName,
		CpuCores: &numCores,
		Mem:      &memoryMB,
		//RootDisk: &rootDiskMB,
	}

	s.Metadata = map[string]string{ // userdata
		metadataKeyIsState:   metadataValueTrue, // bootstrap
		metadataKeyCloudInit: string(userData),
	}
	s.Addresses = []network.Address{{
		Value: "10.0.0.1",
		Type:  network.IPv4Address,
		Scope: network.ScopeCloudLocal,
	}}
	s.Instance = s.NewInstance(c, "spam")
	s.RawInstance = s.Instance.raw
	s.InstName = s.Prefix + "machine-spam"

	s.StartInstArgs = environs.StartInstanceParams{
		InstanceConfig: instanceConfig,
		Tools:          tools,
		Constraints:    cons,
		//Placement: "",
		//DistributionGroup: nil,
	}

	//s.InstanceType = allInstanceTypes[0]
}

func (s *BaseSuiteUnpatched) initNet(c *gc.C) {
	s.Ports = []network.PortRange{{
		FromPort: 80,
		ToPort:   80,
		Protocol: "tcp",
	}}
}

func (s *BaseSuiteUnpatched) setConfig(c *gc.C, cfg *config.Config) {
	s.Config = cfg
	ecfg, err := newValidConfig(cfg, configDefaults)
	c.Assert(err, jc.ErrorIsNil)
	s.EnvConfig = ecfg
	uuid, _ := cfg.UUID()
	s.Env.uuid = uuid
	s.Env.ecfg = s.EnvConfig
	s.Prefix = "juju-" + uuid + "-"
}

func (s *BaseSuiteUnpatched) NewConfig(c *gc.C, updates testing.Attrs) *config.Config {
	if updates == nil {
		updates = make(testing.Attrs)
	}
	var err error
	cfg := testing.EnvironConfig(c)
	cfg, err = cfg.Apply(ConfigAttrs)
	c.Assert(err, jc.ErrorIsNil)
	if raw := updates[cfgNamespace]; raw == nil || raw.(string) == "" {
		updates[cfgNamespace] = cfg.Name()
	}
	cfg, err = cfg.Apply(updates)
	c.Assert(err, jc.ErrorIsNil)
	return cfg
}

func (s *BaseSuiteUnpatched) UpdateConfig(c *gc.C, attrs map[string]interface{}) {
	cfg, err := s.Config.Apply(attrs)
	c.Assert(err, jc.ErrorIsNil)
	s.setConfig(c, cfg)
}

func (s *BaseSuiteUnpatched) NewRawInstance(c *gc.C, name string) *lxdclient.Instance {
	summary := lxdclient.InstanceSummary{
		Name:      name,
		Status:    lxdclient.StatusRunning,
		Hardware:  *s.Hardware,
		Metadata:  s.Metadata,
		Addresses: s.Addresses,
	}
	instanceSpec := lxdclient.InstanceSpec{
		Name:      name,
		Profiles:  []string{},
		Ephemeral: false,
		Metadata:  s.Metadata,
	}
	return lxdclient.NewInstance(summary, &instanceSpec)
}

func (s *BaseSuiteUnpatched) NewInstance(c *gc.C, name string) *environInstance {
	raw := s.NewRawInstance(c, name)
	return newInstance(raw, s.Env)
}

type BaseSuite struct {
	BaseSuiteUnpatched

	Stub       *gitjujutesting.Stub
	Client     *stubClient
	Firewaller *stubFirewaller
	Common     *stubCommon
	Policy     *stubPolicy
}

func (s *BaseSuite) SetUpTest(c *gc.C) {
	s.BaseSuiteUnpatched.SetUpTest(c)

	s.Stub = &gitjujutesting.Stub{}
	s.Client = &stubClient{stub: s.Stub}
	s.Firewaller = &stubFirewaller{stub: s.Stub}
	s.Common = &stubCommon{stub: s.Stub}
	s.Policy = &stubPolicy{stub: s.Stub}

	// Patch out all expensive external deps.
	s.Env.raw = &rawProvider{
		lxdInstances:   s.Client,
		Firewaller:     s.Firewaller,
		policyProvider: s.Policy,
	}
	s.Env.base = s.Common
	//s.PatchValue(&newConnection, func(*environConfig) (gceConnection, error) {
	//	return s.StubConn, nil
	//})
	//s.PatchValue(&supportedArchitectures, s.StubCommon.SupportedArchitectures)
	//s.PatchValue(&bootstrap, s.StubCommon.Bootstrap)
	//s.PatchValue(&destroyEnv, s.StubCommon.Destroy)
	//s.PatchValue(&availabilityZoneAllocations, s.StubCommon.AvailabilityZoneAllocations)
	//s.PatchValue(&buildInstanceSpec, s.StubEnviron.BuildInstanceSpec)
	//s.PatchValue(&getHardwareCharacteristics, s.StubEnviron.GetHardwareCharacteristics)
	//s.PatchValue(&newRawInstance, s.StubEnviron.NewRawInstance)
	//s.PatchValue(&findInstanceSpec, s.StubEnviron.FindInstanceSpec)
	//s.PatchValue(&getInstances, s.StubEnviron.GetInstances)
	//s.PatchValue(&imageMetadataFetch, s.StubImages.ImageMetadataFetch)
}

func (s *BaseSuite) CheckNoAPI(c *gc.C) {
	s.Stub.CheckCalls(c, nil)
}

type stubCommon struct {
	stub *gitjujutesting.Stub

	BootstrapResult *environs.BootstrapResult
}

func (sc *stubCommon) BootstrapEnv(ctx environs.BootstrapContext, params environs.BootstrapParams) (*environs.BootstrapResult, error) {
	sc.stub.AddCall("Bootstrap", ctx, params)
	if err := sc.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return sc.BootstrapResult, nil
}

func (sc *stubCommon) DestroyEnv() error {
	sc.stub.AddCall("Destroy")
	if err := sc.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type stubPolicy struct {
	stub *gitjujutesting.Stub

	Arches []string
}

func (s *stubPolicy) SupportedArchitectures() ([]string, error) {
	s.stub.AddCall("SupportedArchitectures")
	if err := s.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return s.Arches, nil
}

type stubClient struct {
	stub *gitjujutesting.Stub

	Insts []lxdclient.Instance
	Inst  *lxdclient.Instance
}

func (conn *stubClient) Instances(prefix string, statuses ...string) ([]lxdclient.Instance, error) {
	conn.stub.AddCall("Instances", prefix, statuses)
	if err := conn.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return conn.Insts, nil
}

func (conn *stubClient) AddInstance(spec lxdclient.InstanceSpec) (*lxdclient.Instance, error) {
	conn.stub.AddCall("AddInstance", spec)
	if err := conn.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return conn.Inst, nil
}

func (conn *stubClient) RemoveInstances(prefix string, ids ...string) error {
	conn.stub.AddCall("RemoveInstances", prefix, ids)
	if err := conn.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

// TODO(ericsnow) Move stubFirewaller to environs/testing or provider/common/testing.

type stubFirewaller struct {
	stub *gitjujutesting.Stub

	PortRanges []network.PortRange
}

func (fw *stubFirewaller) Ports(fwname string) ([]network.PortRange, error) {
	fw.stub.AddCall("Ports", fwname)
	if err := fw.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return fw.PortRanges, nil
}

func (fw *stubFirewaller) OpenPorts(fwname string, ports ...network.PortRange) error {
	fw.stub.AddCall("OpenPorts", fwname, ports)
	if err := fw.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (fw *stubFirewaller) ClosePorts(fwname string, ports ...network.PortRange) error {
	fw.stub.AddCall("ClosePorts", fwname, ports)
	if err := fw.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
