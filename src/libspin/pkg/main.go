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

// Package pkg provides package management implementations.
//
// Package management
//
// This library abstracts the details of different package managers in order
// to allow the tool to remain agnostic. Each package manager implementation
// follows a specific contract through the Manager interface.
//
// Usage
//
// Using a package manager instance is as simple as gaining a reference to it,
// via the NewPackageManager method.
//      manager, err := pkg.NewManager(pkg.PackageManagerEopkg)
//
// Placement in the lifecycle
//
// The package manager instance is only called upon after the initial root filesystem
// has been initialised by the Builder. During this time, every method of the API
// will be called.
// Implementations should not expect longevity, as soon as the operation set has
// been completed, the rootfs is then finalized ready for packing in whichever form
// is most appropriate to the builder.
package pkg

import (
	"errors"
)

// PackageManager is the type of pkg.Manager that is to be used for packaging
// operations. This type is built-in as there is only a limited set of supported
// package managers at any time.
type PackageManager string

const (
	// PackageManagerEopkg is the package manager used within Solus
	PackageManagerEopkg PackageManager = "eopkg"
)

var (
	// ErrNotYetImplemented is a placeholder until eopkg implementation is done
	ErrNotYetImplemented = errors.New("Not yet implemented!")

	// ErrNotEnoughOps should never, ever happen. So check for it. >_>
	ErrNotEnoughOps = errors.New("Internal error: 0 args passed to ApplyOperations")

	// ErrUnknownOperation is returned when we don't know how to handle an operation
	ErrUnknownOperation = errors.New("Unknown or unsupported operation requested")
)
