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

package pkg

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DeviceNode represents a /dev/ node to be created in chroots
type DeviceNode struct {
	Mode  string // Mode to create the device node in
	Major uint32 // Major ID
	Minor uint32 // Minor ID
	Path  string // Path within a chroot (no / prefix)
}

const (
	// PackageManagerEopkg is the package manager used within Solus
	PackageManagerEopkg = "eopkg"
)

var (
	// ErrNotYetImplemented is a placeholder until eopkg implementation is done
	ErrNotYetImplemented = errors.New("Not yet implemented!")

	// ErrNotEnoughOps should never, ever happen. So check for it. >_>
	ErrNotEnoughOps = errors.New("Internal error: 0 args passed to ApplyOperations")

	// ErrUnknownOperation is returned when we don't know how to handle an operation
	ErrUnknownOperation = errors.New("Unknown or unsupported operation requested")

	// DevNodeRandom is /dev/random
	DevNodeRandom *DeviceNode

	// DevNodeURandom is /dev/urandom
	DevNodeURandom *DeviceNode
)

func init() {
	DevNodeURandom = &DeviceNode{Mode: "00666", Major: 1, Minor: 9, Path: "dev/urandom"}
	DevNodeRandom = &DeviceNode{Mode: "00666", Major: 1, Minor: 8, Path: "dev/random"}
}

// ExecStdoutArgs is a convenience function to execute a command on stdout with
// the given arguments
func ExecStdoutArgs(command string, args []string) error {
	var err error
	// Search the path if necessary
	if !strings.Contains(command, "/") {
		command, err = exec.LookPath(command)
		if err != nil {
			return err
		}
	}
	c := exec.Command(command, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// ChrootExec will run a given command in the chroot directory
func ChrootExec(dir, command string) error {
	cmdArgs := []string{dir, "/bin/sh", "-c", command}
	return ExecStdoutArgs("chroot", cmdArgs)
}

// AddGroup will chroot into the given root and add a group
func AddGroup(root, groupName string, groupID int) error {
	cmd := fmt.Sprintf("groupadd -g %d \"%s\"", groupID, groupName)
	return ChrootExec(root, cmd)
}

// AddUser will chroot into the given root and add a user
func AddUser(root, userName, gecos, home, shell string, uid, gid int) error {
	cmd := fmt.Sprintf("useradd -m -d \"%s\" -s \"%s\" -u %d -g %d \"%s\" -c \"%s\"",
		home, shell, uid, gid, userName, gecos)
	return ChrootExec(root, cmd)
}

// AddSystemUser will chroot into the given root and add a system user
func AddSystemUser(root, userName, gecos, home, shell string, uid, gid int) error {
	cmd := fmt.Sprintf("useradd -m -d \"%s\" -r -s \"%s\" -u %d -g %d \"%s\" -c \"%s\"",
		home, shell, uid, gid, userName, gecos)
	return ChrootExec(root, cmd)
}

// CreateDeviceNode will create the essential nodes in a chroot path
func CreateDeviceNode(root string, node *DeviceNode) error {
	fpath := filepath.Join(root, node.Path)
	cmd := []string{"-m", node.Mode, fpath, "c", fmt.Sprintf("%d", node.Major), fmt.Sprintf("%d", node.Minor)}

	return ExecStdoutArgs("mknod", cmd)
}
