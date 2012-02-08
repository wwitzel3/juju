package charm_test

import (
	"bytes"
	"io"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"launchpad.net/goyaml"
	"launchpad.net/juju/go/charm"
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

func init() {
	// Bazaar can't hold subtle mode differences, so we enforce
	// them here to run more interesting checks below.
	modes := []struct{ path string; mode os.FileMode }{
		{"hooks/install", 0751},
		{"empty", 0750},
		{"src/hello.c", 0614},
	}
	for _, m := range modes {
		err := os.Chmod(filepath.Join(repoDir("dummy"), m.path), m.mode)
		if err != nil {
			panic(err)
		}
	}
}

func checkDummy(c *C, f charm.Charm, path string) {
	c.Assert(f.Revision(), Equals, 1)
	c.Assert(f.Meta().Name, Equals, "dummy")
	c.Assert(f.Config().Options["title"].Default, Equals, "My Title")
	switch f := f.(type) {
	case *charm.Bundle:
		c.Assert(f.Path, Equals, path)

	case *charm.Dir:
		c.Assert(f.Path, Equals, path)
		if path == repoDir("dummy") {
			break // Don't test original Bazaar content.
		}

		info, err := os.Stat(filepath.Join(path, "src", "hello.c"))
		c.Assert(err, IsNil)
		c.Assert(info.Mode() & 0777, Equals, os.FileMode(0644))
		c.Assert(info.Mode() & os.ModeType, Equals, os.FileMode(0))

		info, err = os.Stat(filepath.Join(path, "hooks", "install"))
		c.Assert(err, IsNil)
		c.Assert(info.Mode() & 0777, Equals, os.FileMode(0755))
		c.Assert(info.Mode() & os.ModeType, Equals, os.FileMode(0))

		info, err = os.Stat(filepath.Join(path, "empty"))
		c.Assert(err, IsNil)
		c.Assert(info.Mode() & 0777, Equals, os.FileMode(0755))

		target, err := os.Readlink(filepath.Join(path, "hooks", "symlink"))
		c.Assert(err, IsNil)
		c.Assert(target, Equals, "../target")
	}
}

type YamlHacker map[interface{}]interface{}

func ReadYaml(r io.Reader) YamlHacker {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	m := make(map[interface{}]interface{})
	err = goyaml.Unmarshal(data, m)
	if err != nil {
		panic(err)
	}
	return YamlHacker(m)
}

func (yh YamlHacker) Reader() io.Reader {
	data, err := goyaml.Marshal(yh)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(data)
}
