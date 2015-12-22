// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state_test

import (
	"io"

	"github.com/juju/errors"
	"github.com/juju/testing"

	"github.com/juju/juju/resource"
	"github.com/juju/juju/resource/state"
)

type stubRawState struct {
	stub *testing.Stub

	ReturnPersistence state.Persistence
	ReturnStorage     state.Storage
}

func (s *stubRawState) Persistence() (state.Persistence, error) {
	s.stub.AddCall("Persistence")
	if err := s.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return s.ReturnPersistence, nil
}

func (s *stubRawState) Storage() (state.Storage, error) {
	s.stub.AddCall("Storage")
	if err := s.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return s.ReturnStorage, nil
}

type stubPersistence struct {
	stub *testing.Stub

	ReturnListResources []resource.Resource
}

func (s *stubPersistence) ListResources(serviceID string) ([]resource.Resource, error) {
	s.stub.AddCall("ListResources", serviceID)
	if err := s.stub.NextErr(); err != nil {
		return nil, errors.Trace(err)
	}

	return s.ReturnListResources, nil
}

type stubStorage struct {
	stub *testing.Stub
}

func (s *stubStorage) Put(hash string, r io.Reader, length int64) error {
	s.stub.AddCall("Put", hash, r, length)
	if err := s.stub.NextErr(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type stubReader struct {
	stub *testing.Stub

	ReturnRead int
}

func (s *stubReader) Read(buf []byte) (int, error) {
	s.stub.AddCall("Read", buf)
	if err := s.stub.NextErr(); err != nil {
		return 0, errors.Trace(err)
	}

	return s.ReturnRead, nil
}
