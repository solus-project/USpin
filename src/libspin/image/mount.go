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
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LenSort is to enable reverse length sorting
type LenSort []string

// Len returns the length of the string slice
func (l LenSort) Len() int {
	return len(l)
}

// Swap will swap two elements by index in the string slice
func (l LenSort) Swap(a, b int) {
	l[a], l[b] = l[b], l[a]
}

// Less determines if a is less than b. This is deliberately negated
func (l LenSort) Less(a, b int) bool {
	return len(l[a]) > len(l[b])
}

const (
	// UmountMaxTries is the maximum number of times to try unmounting before
	// resorting to lazy detaches
	UmountMaxTries = 3

	// UmountRetryTime is the length of time to wait in between umounts
	UmountRetryTime = 500 * time.Millisecond
)

// A MountEntry is tracked by the MountManager to enable proper cleanup takes
// place
type MountEntry struct {
	SourcePath string // The source of the mount
	MountPoint string // The destination mount point
}

// Umount will attempt to unmount the given path
func (m *MountEntry) Umount() error {
	return ExecStdout(fmt.Sprintf("umount %v", m.MountPoint))
}

// UmountForce will attempt to forcibly detach the mountpoint
func (m *MountEntry) UmountForce() error {
	return ExecStdout(fmt.Sprintf("umount -f %v", m.MountPoint))
}

// UmountLazy will attempt a lazy detach of the node
func (m *MountEntry) UmountLazy() error {
	return ExecStdout(fmt.Sprintf("umount -l %v", m.MountPoint))
}

// UmountSync will attempt everything possible to umount itself
func (m *MountEntry) UmountSync() error {
	for i := 0; i < UmountMaxTries; i++ {
		if err := m.Umount(); err == nil {
			return nil
		}
		time.Sleep(UmountRetryTime)
	}
	// Still didn't manage to umount it
	if err := m.UmountForce(); err == nil {
		return nil
	}
	return m.UmountLazy()
}

// A MountManager is used to mount and unmount filesystems, and to track them
// so that they are all properly torn down
type MountManager struct {
	mounts map[string]*MountEntry
}

var mountManager *MountManager

func init() {
	mountManager = &MountManager{}
	mountManager.mounts = make(map[string]*MountEntry)
}

// GetMountManager will return the global mount manager
func GetMountManager() *MountManager {
	return mountManager
}

// insertMount will store the given mount point in order to permit deletion of it later
func (m *MountManager) insertMount(sourcepath, destpath string) {
	me := &MountEntry{
		SourcePath: sourcepath,
		MountPoint: destpath,
	}
	m.mounts[destpath] = me
}

// Mount will attempt to mount the given sourcepath at the destpath
func (m *MountManager) Mount(sourcepath, destpath, filesystem string, options ...string) error {
	// Only store the absolute path for the mountpoint
	dpath, err := filepath.Abs(destpath)
	if err != nil {
		return err
	}

	if _, ok := m.mounts[dpath]; ok {
		return fmt.Errorf("Path already known to MountManager: %v", dpath)
	}

	cmd := fmt.Sprintf("mount -t %v %v %v", filesystem, sourcepath, destpath)
	if len(options) > 1 {
		optString := strings.Join(options, ",")
		cmd += fmt.Sprintf(" -o %v", optString)
	}

	if err := ExecStdout(cmd); err != nil {
		return err
	}
	m.insertMount(sourcepath, dpath)
	return nil
}

// BindMount will attempt to mount the given sourcepath at the destpath with a binding
func (m *MountManager) BindMount(sourcepath, destpath, filesystem string) error {
	return m.Mount(sourcepath, destpath, filesystem, "bind")
}

// Unmount will attempt to unmount the given path
func (m *MountManager) Unmount(mountpoint string) error {
	dpath, err := filepath.Abs(mountpoint)
	if err != nil {
		return err
	}
	me, ok := m.mounts[dpath]
	if !ok {
		return fmt.Errorf("Attempting to umount unknown path to manager: %v", dpath)
	}
	err = me.UmountSync()
	delete(m.mounts, dpath)
	return err
}

// UnmountAll will attempt to unmount all registered mountpoints
func (m *MountManager) UnmountAll() {
	ExecStdout("sync")
	var keys []string
	for key := range m.mounts {
		keys = append(keys, key)
	}
	sort.Sort(LenSort(keys))
	for _, key := range keys {
		if err := m.Unmount(key); err != nil {
			log.Error(err)
		}
	}
}
