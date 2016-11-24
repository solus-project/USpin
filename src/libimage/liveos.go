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
	img       *libspin.ImageSpec
	rootfsImg string
	rootfsDir string
	deployDir string
	workspace string
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
	l.rootfsImg = l.JoinPath("rootfs.img")
	l.rootfsDir = l.JoinPath("rootfs")
	l.deployDir = l.JoinPath("deploy")

	// As and when we add new directories, populate them here
	requiredDirs := []string{
		l.workspace,
		l.rootfsDir,
		l.deployDir,
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
	rootSize := l.img.Config.LiveOS.RootfsSize
	rootFormat := l.img.Config.LiveOS.RootfsFormat
	if err := CreateSparseFile(l.rootfsImg, rootSize); err != nil {
		return err
	}
	if err := FormatAs(l.rootfsImg, rootFormat); err != nil {
		return err
	}
	return nil
}

// Cleanup currently does nothing within this builder
func (l *LiveOSBuilder) Cleanup() {
	log.Info("Cleaning up")
}
