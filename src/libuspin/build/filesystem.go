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

package build

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
)

// FilesystemFormatFunc is the prototype for functions that format filesystems
// to ensure we can use dedicated functions that can handle filesystem paths
// correctly (i.e. spaces)
type FilesystemFormatFunc func(filename string) error

// A FilesystemCheckFunc is a function prototype for performing filesystem
// checks, i.e. a rootfs.img after unmounting
type FilesystemCheckFunc func(filename string) error

var filesystemCommands map[string]FilesystemFormatFunc
var checkCommands map[string]FilesystemCheckFunc

func formatExt4(filename string) error {
	// Format it
	if err := ExecStdoutArgs("mkfs", []string{"-t", "ext4", "-F", filename}); err != nil {
		return err
	}
	// Set the mount count so it doesn't get fsck'd during live boot
	return ExecStdoutArgs("tune2fs", []string{"-c0", "-i0", filename})
}

func checkExt4(filename string) error {
	// Check it for errors
	if err := ExecStdoutArgs("e2fsck", []string{"-y", filename}); err != nil {
		return err
	}
	// Force fix any issues now
	return ExecStdoutArgs("e2fsck", []string{"-y", "-f", filename})
}

func init() {
	// Initialise the command maps
	filesystemCommands = make(map[string]FilesystemFormatFunc)
	checkCommands = make(map[string]FilesystemCheckFunc)

	filesystemCommands["ext4"] = formatExt4
	checkCommands["ext4"] = checkExt4
}

// FormatAs will format the given path with the filesystem specified.
// Note: You should only use this with image paths, it's dangerous!
func FormatAs(filename, filesystem string) error {
	command, ok := filesystemCommands[filesystem]
	if !ok {
		return fmt.Errorf("Cannot format with unknown filesystem '%v'", filesystem)
	}
	log.WithFields(log.Fields{
		"filename":   filename,
		"filesystem": filesystem,
	}).Info("Formatting filesystem")
	return command(filename)
}

// CheckFS will try to check/fix the filesystems pointed to by filename
// using the helpers denoted by filesystem.
// This should only be used for internal image code on loopback devices!
func CheckFS(filename, filesystem string) error {
	command, ok := checkCommands[filesystem]
	if !ok {
		return fmt.Errorf("Cannot check with unknown filesystem '%v'", filesystem)
	}
	log.WithFields(log.Fields{
		"filename":   filename,
		"filesystem": filesystem,
	}).Info("Checking filesystem")
	return command(filename)
}
