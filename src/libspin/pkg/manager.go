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

package pkg

import (
	"errors"
	"libspin/config"
	"libspin/spec"
)

// Manager is the interface that should be implemented by vendors to enable
// solspin to understand them and construct images according to their particulars.
type Manager interface {

	// Init will allwo implementations to initialise themselves according to any
	// particulars in the image configuration, allowing it to make better informed
	// decisions later on.
	Init(conf *config.ImageConfiguration) error

	// InitRoot implementations should set up the root filesystem to handle any
	// quirks prior to installing packages. This also allows manipulating the
	// filesystem layout, i.e. for usr-merge situations, or for working around
	// default directories created by a host-side package manager tool.
	InitRoot(root string) error

	// ApplyOperations should apply all of the given operations in bulk, as they
	// are always guaranteed to have the same type.
	ApplyOperations(ops []spec.Operation) error

	// Cleanup may be called at any time, and the package manager implementation
	// should ensure it cleans anything it did in the past, such as closing open
	// processes.
	Cleanup() error
}

// NewManager will return an appropriate package manager instance for
// the given name, if it exists.
func NewManager(name string) (Manager, error) {
	switch name {
	case PackageManagerEopkg:
		return NewEopkgManager(), nil
	default:
		return nil, errors.New("Not yet implemented")
	}
}
