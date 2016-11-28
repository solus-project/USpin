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
}

// A LoaderType is a pseudo enum type for the bootloader to restrict to
// supported implementations
type LoaderType string

const (
	// LoaderTypeSyslinux refers to syslinux + isolinux
	LoaderTypeSyslinux LoaderType = "syslinux"
)

var (
	// ErrNotYetImplemented is just a placeholder
	ErrNotYetImplemented = errors.New("Not yet implemented")
)

// NewLoader will create a new Loader instance for the given name, if supported
func NewLoader(impl LoaderType) (Loader, error) {
	switch impl {
	case LoaderTypeSyslinux:
		return NewSyslinuxLoader(), nil
	default:
		return nil, ErrNotYetImplemented
	}
}
