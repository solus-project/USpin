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
	"libspin/config"
	"libspin/spec"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	// ErrNotYetImplemented is a placeholder until eopkg implementation is done
	ErrNotYetImplemented = errors.New("Not yet implemented!")

	// ErrNotEnoughOps should never, ever happen. So check for it. >_>
	ErrNotEnoughOps = errors.New("Internal error: 0 args passed to ApplyOperations")

	// ErrUnknownOperation is returned when we don't know how to handle an operation
	ErrUnknownOperation = errors.New("Unknown or unsupported operation requested")
)

// EopkgManager is used to apply operations with the eopkg package manager
// for Solus systems.
type EopkgManager struct {
	root string // rootfs path
}

// NewEopkgManager will return a newly initialised EopkgManager
func NewEopkgManager() *EopkgManager {
	return &EopkgManager{}
}

// Init will attempt to initialise EopkgManager from the given configuration
func (e *EopkgManager) Init(conf *config.ImageConfiguration) error {
	// Ensure the system has eopkg available first!
	if _, err := exec.LookPath("eopkg"); err != nil {
		return err
	}
	return nil
}

// InitRoot will set up the filesystem root in accordance with eopkg needs
func (e *EopkgManager) InitRoot(root string) error {
	e.root = root

	// Ensures we don't end up with /var/lock vs /run/lock nonsense
	reqDirs := []string{
		"run/lock",
		"var",
	}

	// Construct the required directories in the tree
	for _, dir := range reqDirs {
		dirPath := filepath.Join(root, dir)
		if err := os.MkdirAll(dirPath, 00755); err != nil {
			return err
		}
	}

	if err := os.Symlink("../run/lock", filepath.Join(root, "var", "lock")); err != nil {
		return err
	}
	if err := os.Symlink("../run", filepath.Join(root, "var", "run")); err != nil {
		return err
	}
	return nil
}

// ApplyOperations will apply the given set of operations via eopkg
func (e *EopkgManager) ApplyOperations(ops []spec.Operation) error {
	if len(ops) == 0 {
		return ErrNotEnoughOps
	}
	switch ops[0].(type) {
	case *spec.OpRepo:
		return e.addRepos(ops)
	default:
		return ErrUnknownOperation
	}
}

// Cleanup will cleanup the rootfs at any given point
func (e *EopkgManager) Cleanup() error {
	return ErrNotYetImplemented
}

// Eopkg specific functions

// Add a repository to the target
func (e *EopkgManager) addRepos(ops []spec.Operation) error {
	for _, repo := range ops {
		r := repo.(*spec.OpRepo)
		cmd := fmt.Sprintf("eopkg -D %s add-repo %s %s", e.root, r.RepoName, r.RepoURI)
		fmt.Fprintf(os.Stderr, "TODO: %s\n", cmd)
	}
	return ErrNotYetImplemented
}
