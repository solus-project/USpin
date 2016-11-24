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
	"errors"
	"libspin"
)

var (
	// ErrNotYetImplemented is just used until we actually implement some code....
	ErrNotYetImplemented = errors.New("Not yet implemented")
)

// A LiveOSBuilder is responsible for building ISO format images that are USB
// compatible. It is the "LiveCD" type of Builder
type LiveOSBuilder struct {
	img *libspin.ImageSpec
}

// NewLiveOSBuilder should only be used by builder.go
func NewLiveOSBuilder() *LiveOSBuilder {
	return &LiveOSBuilder{}
}

// Init will initialise a LiveOSBuilder from the given spec
func (l *LiveOSBuilder) Init(img *libspin.ImageSpec) error {
	l.img = img
	return nil
}

// PrepareWorkspace sets up the required directories for the LiveOSBuilder
func (l *LiveOSBuilder) PrepareWorkspace() error {
	return ErrNotYetImplemented
}
