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
	"strings"
	"syscall"
)

// A MountManager is used to mount and unmount filesystems, and to track them
// so that they are all properly torn down
type MountManager struct {
}

var mountManager *MountManager

func init() {
	mountManager = &MountManager{}
}

// GetMountManager will return the global mount manager
func GetMountManager() *MountManager {
	return mountManager
}

// MountPath will attempt to mount the given sourcepath at the destpath
func (m *MountManager) MountPath(sourcepath, destpath, filesystem string, flags uintptr, options ...string) error {
	optString := ""
	if len(options) > 1 {
		optString = strings.Join(options, ",")
	}
	er := syscall.Mount(sourcepath, destpath, filesystem, flags, optString)
	return er
}

// BindMount will attempt to mount the given sourcepath at the destpath with a binding
func (m *MountManager) BindMount(sourcepath, destpath, filesystem string, options ...string) error {
	return m.MountPath(sourcepath, destpath, filesystem, syscall.MS_BIND, options...)
}

// Mount will attempt to mount the given sourcepath at the destpath with default options
func (m *MountManager) Mount(sourcepath, destpath, filesystem string, options ...string) error {
	return m.MountPath(sourcepath, destpath, filesystem, 0, options...)
}
