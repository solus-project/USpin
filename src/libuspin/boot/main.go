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
)

// A Loader provides abstraction around various bootloader implementations.
type Loader interface {
	Init() error
	// GetCapabilities returns the supported capabilities of this bootloader implementation
	GetCapabilities() Capability
}

// A LoaderType is a pseudo enum type for the bootloader to restrict to
// supported implementations
type LoaderType string

// Capability refers to the type of operations that a bootloader supports
type Capability uint8

const (
	// CapInstallUEFI means the bootloader supports UEFI loading
	CapInstallUEFI Capability = 1 << 1

	// CapInstallLegacy means the bootloader supports MBR/legacy loading
	CapInstallLegacy Capability = 1 << 2

	// CapInstallISO is used for bootloaders reporting ISO support
	CapInstallISO Capability = 1 << 3

	// CapInstallRaw is reported by bootloaders that can install to block devices
	CapInstallRaw Capability = 1 << 4
)

const (
	// LoaderTypeSyslinux refers to syslinux + isolinux
	LoaderTypeSyslinux LoaderType = "syslinux"
)

var (
	// ErrNotYetImplemented is just a placeholder
	ErrNotYetImplemented = errors.New("Not yet implemented")

	// ErrUnknownLoader is reported for an unknown bootloader
	ErrUnknownLoader = errors.New("Unknown bootloader configured")
)

// NewLoader will create a new Loader instance for the given name, if supported
func NewLoader(impl LoaderType) (Loader, error) {
	switch impl {
	case LoaderTypeSyslinux:
		return NewSyslinuxLoader(), nil
	default:
		return nil, ErrUnknownLoader
	}
}

// InitLoaders will attempt to return an initialised set of loaders as a helper
// to other Builder implementations
func InitLoaders(loaderType []LoaderType) ([]Loader, error) {
	var ret []Loader

	for _, name := range loaderType {
		if loader, err := NewLoader(name); err == nil {
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
