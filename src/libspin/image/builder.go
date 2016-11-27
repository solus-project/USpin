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

package image

import (
	"fmt"
	"libspin"
	"libspin/config"
)

// A Builder is the contract definition for all image builders, and the implementations
// will worry about the exact particulars.
type Builder interface {

	// Init will initialise the builder from the given configuration
	Init(img *libspin.ImageSpec) error

	// PrepareWorkspace will attempt to do all prework for setting up the image
	// deployment areas and such.
	PrepareWorkspace() error

	// CreateStorage is used by implementations to create any initial backing storage
	// they will require, i.e. the place where we install packages to. No processes
	// should be spawned within it, nor should it be mounted, at this point.
	CreateStorage() error

	// MountStorage should be used by the implementation if it needs to do any mounting
	// to allow the package manager instance to start installing packages and such into
	// the rootfs.
	MountStorage() error

	// UnmountStorage should be used by the implementation to tear down any mounts
	// previously erected in MountStorage() for package manager operations, prior
	// to image finalisation
	UnmountStorage() error

	// FinalizeImage will complete the stage2 part of the image construction, whereby
	// all manual input from package managers, etc, is no longer needed.
	FinalizeImage() error

	// GetRootDir is used by implementations to return the root directory for the
	// OS files
	GetRootDir() string

	// Cleanup should be used by implementations to do any required cleanup operations,
	// including killing processes, unmounting anything, etc.
	Cleanup()
}

// NewBuilder will try to return a builder for the given type
func NewBuilder(name config.ImageType) (Builder, error) {
	switch name {
	case config.ImageTypeLiveOS:
		return NewLiveOSBuilder(), nil
	default:
		return nil, fmt.Errorf("Unknown builder: %v", name)
	}
}
