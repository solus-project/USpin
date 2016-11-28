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

// A Bootloader provides abstraction around various bootloader implementations.
type Bootloader interface {
	Init() error
}

// A Type is a pseudo enum type for the bootloader to restrict to
// supported implementations
type Type string

const (
	// BootloaderTypeSyslinux refers to syslinux + isolinux
	BootloaderTypeSyslinux Type = "syslinux"
)

// New will create a new Bootloader instance for the given name, if supported
func New(impl Type) (Bootloader, error) {
	return nil, errors.New("Not yet implemented")
}
