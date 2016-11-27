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
	"errors"
	"libspin"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	// ErrNotYetImplemented is just used until we actually implement some code....
	ErrNotYetImplemented = errors.New("Not yet implemented")

	requiredBinaries []string
)

const (
	// DefaultImageSize is the size of the rootfs we try to create (4GB)
	DefaultImageSize = 4000
)

func init() {
	requiredBinaries = []string{
		"isohybrid",
		"mksquashfs",
		"xorriso",
	}
}

// A LiveOSBuilder is responsible for building ISO format images that are USB
// compatible. It is the "LiveCD" type of Builder
type LiveOSBuilder struct {
	img            *libspin.ImageSpec
	rootfsImg      string
	rootfsDir      string
	rootfsFormat   string
	rootfsSize     int
	deployDir      string
	liveosDir      string
	liveStagingDir string
	workspace      string
}

// NewLiveOSBuilder should only be used by builder.go
func NewLiveOSBuilder() *LiveOSBuilder {
	return &LiveOSBuilder{}
}

// Init will initialise a LiveOSBuilder from the given spec
func (l *LiveOSBuilder) Init(img *libspin.ImageSpec) error {
	l.img = img

	// Ensure all required binaries are available before we go doing anything.
	for _, bin := range requiredBinaries {
		if _, err := exec.LookPath(bin); err != nil {
			return err
		}
	}

	// rootfs.img particulars
	l.rootfsFormat = l.img.Config.LiveOS.RootfsFormat
	l.rootfsSize = l.img.Config.LiveOS.RootfsSize

	return nil
}

// JoinPath is a helper to join paths onto our root workspace directory
func (l *LiveOSBuilder) JoinPath(paths ...string) string {
	return filepath.Join(l.workspace, filepath.Join(paths...))
}

// PrepareWorkspace sets up the required directories for the LiveOSBuilder
func (l *LiveOSBuilder) PrepareWorkspace() error {
	var err error
	if l.workspace, err = filepath.Abs("./workspace"); err != nil {
		return err
	}
	// Initialise our base variables
	l.rootfsDir = l.JoinPath("rootfs")
	l.deployDir = l.JoinPath("deploy")
	// Inside the ISO target
	l.liveosDir = l.JoinPath("deploy", "LiveOS")
	// Inside the workspace only
	l.liveStagingDir = l.JoinPath("LiveOS")
	l.rootfsImg = l.JoinPath("LiveOS", "rootfs.img")

	// As and when we add new directories, populate them here
	requiredDirs := []string{
		l.workspace,
		l.rootfsDir,
		l.deployDir,
		l.liveosDir,
		l.liveStagingDir,
	}

	// Create all required directories
	for _, dir := range requiredDirs {
		if err = os.MkdirAll(dir, 00755); err != nil {
			return err
		}
	}

	return nil
}

// CreateStorage will create the rootfs.img in which we will contain the
// Live OS
func (l *LiveOSBuilder) CreateStorage() error {
	if err := CreateSparseFile(l.rootfsImg, l.rootfsSize); err != nil {
		return err
	}
	if err := FormatAs(l.rootfsImg, l.rootfsFormat); err != nil {
		return err
	}
	return nil
}

// Cleanup currently does nothing within this builder
func (l *LiveOSBuilder) Cleanup() {
	log.Info("Cleaning up")
	GetMountManager().UnmountAll()
}

// MountStorage will mount the rootfs.img so that the package manager can
// take over
func (l *LiveOSBuilder) MountStorage() error {
	return GetMountManager().Mount(l.rootfsImg, l.rootfsDir, l.rootfsFormat, "loop")
}

// UnmountStorage will unmount the rootfs.img from earlier
func (l *LiveOSBuilder) UnmountStorage() error {
	return GetMountManager().Unmount(l.rootfsDir)
}

// GetRootDir returns the path to the mounted rootfs.img
func (l *LiveOSBuilder) GetRootDir() string {
	return l.rootfsDir
}

// FinalizeImage will go ahead and finish up the ISO construction
func (l *LiveOSBuilder) FinalizeImage() error {
	return ErrNotYetImplemented
}
