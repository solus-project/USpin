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

package libspin

import (
	"errors"
	"libspin/config"
	"libspin/spec"
)

var (
	// ErrNotYetImplemented is a placeholder until eopkg implementation is done
	ErrNotYetImplemented = errors.New("Not yet implemented!")
)

// EopkgManager is used to apply operations with the eopkg package manager
// for Solus systems.
type EopkgManager struct{}

// NewEopkgManager will return a newly initialised EopkgManager
func NewEopkgManager() *EopkgManager {
	return &EopkgManager{}
}

// Init will attempt to initialise EopkgManager from the given configuration
func (e *EopkgManager) Init(conf *config.ImageConfiguration) error {
	return ErrNotYetImplemented
}

// InitRoot will set up the filesystem root in accordance with eopkg needs
func (e *EopkgManager) InitRoot() error {
	return ErrNotYetImplemented
}

// ApplyOperations will apply the given set of operations via eopkg
func (e *EopkgManager) ApplyOperations(ops []spec.Operation) error {
	return ErrNotYetImplemented
}

// Cleanup will cleanup the rootfs at any given point
func (e *EopkgManager) Cleanup() error {
	return ErrNotYetImplemented
}
