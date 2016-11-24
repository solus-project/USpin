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

package libimage

import (
	"fmt"
	"libspin"
)

// A Builder is the contract definition for all image builders, and the implementations
// will worry about the exact particulars.
type Builder interface {

	// Init will initialise the builder from the given configuration
	Init(img *libspin.ImageSpec) error

	// PrepareWorkspace will attempt to do all prework for setting up the image
	// deployment areas and such.
	PrepareWorkspace() error
}

// NewBuilder will try to return a builder for the given type
func NewBuilder(name string) (*Builder, error) {
	switch name {
	default:
		{
			return nil, fmt.Errorf("Unknown builder: %v", name)
		}
	}
}
