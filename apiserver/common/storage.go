// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package common

import (
	"path"

	"github.com/juju/errors"
	"github.com/juju/names"

	"github.com/juju/juju/state"
	"github.com/juju/juju/storage"
)

// StorageInterface is an interface for obtaining information about storage
// instances and related entities.
type StorageInterface interface {
	// StorageInstance returns the state.StorageInstance corresponding
	// to the specified storage tag.
	StorageInstance(names.StorageTag) (state.StorageInstance, error)

	// StorageInstanceFilesystem returns the state.Filesystem assigned
	// to the storage instance with the specified storage tag.
	StorageInstanceFilesystem(names.StorageTag) (state.Filesystem, error)

	// StorageInstanceVolume returns the state.Volume assigned to the
	// storage instance with the specified storage tag.
	StorageInstanceVolume(names.StorageTag) (state.Volume, error)

	// FilesystemAttachment returns the state.FilesystemAttachment
	// corresponding to the identified machine and filesystem.
	FilesystemAttachment(names.MachineTag, names.FilesystemTag) (state.FilesystemAttachment, error)

	// VolumeAttachment returns the state.VolumeAttachment corresponding
	// to the identified machine and volume.
	VolumeAttachment(names.MachineTag, names.VolumeTag) (state.VolumeAttachment, error)

	// WatchFilesystemAttachment watches for changes to the filesystem
	// attachment corresponding to the identfified machien and filesystem.
	WatchFilesystemAttachment(names.MachineTag, names.FilesystemTag) state.NotifyWatcher

	// WatchVolumeAttachment watches for changes to the volume attachment
	// corresponding to the identfified machien and volume.
	WatchVolumeAttachment(names.MachineTag, names.VolumeTag) state.NotifyWatcher
}

// StorageAttachmentInfo returns the StorageAttachmentInfo for the specified
// StorageAttachment by gathering information from related entities (volumes,
// filesystems).
func StorageAttachmentInfo(
	st StorageInterface,
	att state.StorageAttachment,
	machineTag names.MachineTag,
) (*storage.StorageAttachmentInfo, error) {
	storageInstance, err := st.StorageInstance(att.StorageInstance())
	if err != nil {
		return nil, errors.Annotate(err, "getting storage instance")
	}
	switch storageInstance.Kind() {
	case state.StorageKindBlock:
		return volumeStorageAttachmentInfo(st, storageInstance, machineTag)
	case state.StorageKindFilesystem:
		return filesystemStorageAttachmentInfo(st, storageInstance, machineTag)
	}
	return nil, errors.Errorf("invalid storage kind %v", storageInstance.Kind())
}

func volumeStorageAttachmentInfo(
	st StorageInterface,
	storageInstance state.StorageInstance,
	machineTag names.MachineTag,
) (*storage.StorageAttachmentInfo, error) {
	storageTag := storageInstance.StorageTag()
	volume, err := st.StorageInstanceVolume(storageTag)
	if err != nil {
		return nil, errors.Annotate(err, "getting volume")
	}
	volumeInfo, err := volume.Info()
	if err != nil {
		return nil, errors.Annotate(err, "getting volume info")
	}
	volumeAttachment, err := st.VolumeAttachment(machineTag, volume.VolumeTag())
	if err != nil {
		return nil, errors.Annotate(err, "getting volume attachment")
	}
	volumeAttachmentInfo, err := volumeAttachment.Info()
	if err != nil {
		return nil, errors.Annotate(err, "getting volume attachment info")
	}
	devicePath, err := volumeAttachmentDevicePath(
		volumeInfo,
		volumeAttachmentInfo,
	)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &storage.StorageAttachmentInfo{
		storage.StorageKindBlock,
		devicePath,
	}, nil
}

func filesystemStorageAttachmentInfo(
	st StorageInterface,
	storageInstance state.StorageInstance,
	machineTag names.MachineTag,
) (*storage.StorageAttachmentInfo, error) {
	storageTag := storageInstance.StorageTag()
	filesystem, err := st.StorageInstanceFilesystem(storageTag)
	if err != nil {
		return nil, errors.Annotate(err, "getting filesystem")
	}
	filesystemAttachment, err := st.FilesystemAttachment(machineTag, filesystem.FilesystemTag())
	if err != nil {
		return nil, errors.Annotate(err, "getting filesystem attachment")
	}
	filesystemAttachmentInfo, err := filesystemAttachment.Info()
	if err != nil {
		return nil, errors.Annotate(err, "getting filesystem attachment info")
	}
	return &storage.StorageAttachmentInfo{
		storage.StorageKindFilesystem,
		filesystemAttachmentInfo.MountPoint,
	}, nil
}

// WatchStorageAttachmentInfo returns a state.NotifyWatcher that reacts to changes
// to the VolumeAttachmentInfo or FilesystemAttachmentInfo corresponding to the tags
// specified.
func WatchStorageAttachmentInfo(
	st StorageInterface,
	storageTag names.StorageTag,
	machineTag names.MachineTag,
) (state.NotifyWatcher, error) {
	storageInstance, err := st.StorageInstance(storageTag)
	if err != nil {
		return nil, errors.Annotate(err, "getting storage instance")
	}
	switch storageInstance.Kind() {
	case state.StorageKindBlock:
		volume, err := st.StorageInstanceVolume(storageTag)
		if err != nil {
			return nil, errors.Annotate(err, "getting storage volume")
		}
		return st.WatchVolumeAttachment(machineTag, volume.VolumeTag()), nil
	case state.StorageKindFilesystem:
		filesystem, err := st.StorageInstanceFilesystem(storageTag)
		if err != nil {
			return nil, errors.Annotate(err, "getting storage filesystem")
		}
		return st.WatchFilesystemAttachment(machineTag, filesystem.FilesystemTag()), nil
	}
	return nil, errors.Errorf("invalid storage kind %v", storageInstance.Kind())
}

var errNoDevicePath = errors.New("cannot determine device path: no serial or persistent device name")

// volumeAttachmentDevicePath returns the absolute device path for
// a volume attachment. The value is only meaningful in the context
// of the machine that the volume is attached to.
func volumeAttachmentDevicePath(
	volumeInfo state.VolumeInfo,
	volumeAttachmentInfo state.VolumeAttachmentInfo,
) (string, error) {
	if volumeInfo.Serial != "" {
		return path.Join("/dev/disk/by-id", volumeInfo.Serial), nil
	} else if volumeAttachmentInfo.DeviceName != "" {
		return path.Join("/dev", volumeAttachmentInfo.DeviceName), nil
	}
	return "", errNoDevicePath
}
