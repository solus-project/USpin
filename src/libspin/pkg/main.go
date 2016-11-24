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
	"os"
	"os/exec"
	"strings"
)

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
	cmdArgs := []string{dir, "/bin/bash", "-c", command}
	return ExecStdoutArgs("chroot", cmdArgs)
}
