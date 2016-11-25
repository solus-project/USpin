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
	"os"
	"os/exec"
	"strings"
)

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
