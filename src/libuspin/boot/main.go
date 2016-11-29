//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package boot provides implementations to help manage the bootloader
// setup and installation on various medium types.
package boot

import (
	"errors"
	"libuspin/config"
)

// A FileType is a named special file
type FileType string

const (
	// FileTypeBootElToritoCatalog should be within the boot directory on an ISO's
	// boot directory, depending on the bootloader implementation. Note this is not
	// a "real" file but one created by xorriso.
	FileTypeBootElToritoCatalog FileType = "boot.cat"

	// FileTypeBootElToritoBinary is the isolinux.bin style file, i.e. the actual
	// bootloader itself.
	FileTypeBootElToritoBinary FileType = "boot.bin"

	// FileTypeBootMBR is the ISO MBR file. This permits hybrid ISO generation for
	// both USB & CD.
	FileTypeBootMBR FileType = "boot.mbr"
)

// A Loader provides abstraction around various bootloader implementations.
type Loader interface {

	// Init is used by bootloader implementations to assert sanity & host-side tooling
	// is present.
	Init(c *config.ImageConfiguration) error

	// GetCapabilities returns the supported capabilities of this bootloader implementation
	GetCapabilities() Capability

	// Install will have this boot loader implementation install using the given boot
	// configuration.
	Install(mode Capability, c ConfigurationSource) error

	// GetSpecialFile is to enable a private communication method between builder and
	// loader to get well known file types, i.e. "boot.cat"
	// These should *always* respect the relative path rules of the builders
	GetSpecialFile(t FileType) string
}

// ConfigurationSource should be implemented by Builder instances (or their helpers)
// to help with bootloader installation.
type ConfigurationSource interface {

	// JoinRootPath is used by implementations to join a resource path on the rootfs
	JoinRootPath(paths ...string) string

	// JoinDeployPath is used by implementations to join a resource path on the deployment
	// directory.
	// This is mostly of interest to ISO deployments
	JoinDeployPath(paths ...string) string

	// GetRootDevice should return the device used for / mount, if relevant. This
	// should return "" unless used in Raw capability
	GetRootDevice() string

	// GetBootDevice should return the device used for the boot mount, if relevant.
	// This should return "" unless used in the Raw capability
	GetBootDevice() string

	// GetKernel should return the default kernel, configured with the correct
	// asset path
	GetKernel() *Kernel
}

// Capability refers to the type of operations that a bootloader supports
type Capability uint8

const (
	// CapInstallUEFI means the bootloader supports UEFI loading
	CapInstallUEFI Capability = 1 << iota

	// CapInstallLegacy means the bootloader supports MBR/legacy loading
	CapInstallLegacy Capability = 1 << iota

	// CapInstallISO is used for bootloaders reporting ISO support
	CapInstallISO Capability = 1 << iota

	// CapInstallRaw is reported by bootloaders that can install to block devices
	CapInstallRaw Capability = 1 << iota
)

var (
	// ErrNotYetImplemented is just a placeholder
	ErrNotYetImplemented = errors.New("Not yet implemented")

	// ErrUnknownLoader is reported for an unknown bootloader
	ErrUnknownLoader = errors.New("Unknown bootloader configured")
)

// NewLoader will create a new Loader instance for the given name, if supported
func NewLoader(impl config.LoaderType) (Loader, error) {
	switch impl {
	case config.LoaderTypeSyslinux:
		return NewSyslinuxLoader(), nil
	default:
		return nil, ErrUnknownLoader
	}
}

// InitLoaders will attempt to return an initialised set of loaders as a helper
// to other Builder implementations
func InitLoaders(c *config.ImageConfiguration, loaderType []config.LoaderType) ([]Loader, error) {
	var ret []Loader

	for _, name := range loaderType {
		if loader, err := NewLoader(name); err == nil {
			// Init the loader
			if err := loader.Init(c); err != nil {
				return nil, err
			}
			ret = append(ret, loader)
		} else {
			return nil, err
		}
	}
	return ret, nil
}

// GetLoaderWithMask will look in the set of loaders for the given mask.
// Note it will always return the *first* one found, i.e. the first one specified
// in the configuration
func GetLoaderWithMask(loaders []Loader, mask Capability) Loader {
	for _, i := range loaders {
		if i.GetCapabilities()&mask == mask {
			return i
		}
	}
	return nil
}

// HaveLoaderWithMask will look in the set of loaders for the given mask, and
// simply return if the given mask is supported or not.
func HaveLoaderWithMask(loaders []Loader, mask Capability) bool {
	return GetLoaderWithMask(loaders, mask) != nil
}
