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
	"io/ioutil"
	"libspin/config"
	"libspin/image"
	"libspin/spec"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	// EopkgCacheDirectory is where we'll bind mount to provide package caching
	// to speed up subsequent image builds.
	// This will be mounted at $rootfs/var/cache/eopkg/packages.
	// It uses the evobuild directory for consistency with evobuild, so that
	// Solus developers only need one cache system wide.
	EopkgCacheDirectory = "/var/lib/evobuild/packages"
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
		// Enables our bind mounting for caching
		"var/cache/eopkg/packages",
	}

	// Construct the required directories in the tree
	for _, dir := range reqDirs {
		dirPath := filepath.Join(root, dir)
		if err := os.MkdirAll(dirPath, 00755); err != nil {
			return err
		}
	}

	// Attempt to create the system wide cache directory
	if err := os.MkdirAll(EopkgCacheDirectory, 00755); err != nil {
		return err
	}

	if err := os.Symlink("../run/lock", filepath.Join(root, "var", "lock")); err != nil {
		return err
	}
	if err := os.Symlink("../run", filepath.Join(root, "var", "run")); err != nil {
		return err
	}

	// Now attempt to bind mount the cache directory to be .. well. usable
	cacheTarget := filepath.Join(root, "var", "cache", "eopkg", "packages")
	if err := image.GetMountManager().BindMount(EopkgCacheDirectory, cacheTarget); err != nil {
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
	case *spec.OpGroup:
		ignoreSafety := ops[0].(*spec.OpGroup).IgnoreSafety
		return e.installComponents(ops, ignoreSafety)
	case *spec.OpPackage:
		ignoreSafety := ops[0].(*spec.OpPackage).IgnoreSafety
		return e.installPackages(ops, ignoreSafety)
	default:
		return ErrUnknownOperation
	}
}

// FinalizeRoot will configure all of the eopkgs installed in the system, and
// ensure that dbus, etc, works.
func (e *EopkgManager) FinalizeRoot() error {
	if err := e.copyBaselayout(); err != nil {
		return err
	}
	if err := e.configureDbus(); err != nil {
		return err
	}
	if err := CreateDeviceNode(e.root, DevNodeRandom); err != nil {
		return err
	}
	if err := CreateDeviceNode(e.root, DevNodeURandom); err != nil {
		return err
	}
	if err := e.startDBUS(); err != nil {
		return err
	}
	if err := ChrootExec(e.root, "eopkg configure-pending"); err != nil {
		e.killDBUS()
		return err
	}
	return e.killDBUS()
}

// This needs to die in a fire and will not be supported when sol replaces eopkg
func (e *EopkgManager) copyBaselayout() error {
	var files []os.FileInfo
	var err error

	// elements of /usr/share/baselayout are copied to /etc/ - ANTI STATELESS
	baseDir := filepath.Join(e.root, "usr", "share", "baselayout")
	tgtDir := filepath.Join(e.root, "etc")
	if files, err = ioutil.ReadDir(baseDir); err != nil {
		return err
	}

	for _, file := range files {
		srcPath := filepath.Join(baseDir, file.Name())
		tgtPath := filepath.Join(tgtDir, file.Name())

		if err = image.CopyFile(srcPath, tgtPath); err != nil {
			return err
		}
	}
	return nil
}

// Attempt to start dbus in the root..
func (e *EopkgManager) startDBUS() error {
	if err := ChrootExec(e.root, "dbus-uuidgen --ensure"); err != nil {
		return err
	}
	if err := ChrootExec(e.root, "dbus-daemon --system"); err != nil {
		return err
	}
	return nil
}

// killDBUS will stop dbus again
// TODO: Remove the file
func (e *EopkgManager) killDBUS() error {
	fpath := filepath.Join(e.root, "var/run/dbus/pid")
	var b []byte
	var err error
	var f *os.File

	if f, err = os.Open(fpath); err != nil {
		return err
	}
	defer f.Close()

	if b, err = ioutil.ReadAll(f); err != nil {
		return err
	}

	pid := strings.Split(string(b), "\n")[0]
	return ExecStdoutArgs("kill", []string{"-9", pid})
}

// This is also largely anti-stateless but is required just to get dbus running
// so we can configure-pending. sol can't come quick enough...
func (e *EopkgManager) configureDbus() error {
	if err := AddGroup(e.root, "messagebus", 18); err != nil {
		return err
	}
	if err := AddSystemUser(e.root, "messagebus", "D-Bus Message Daemon", "/var/run/dbus", "/bin/false", 18, 18); err != nil {
		return err
	}
	return nil
}

// Cleanup will cleanup the rootfs at any given point
func (e *EopkgManager) Cleanup() error {
	return ErrNotYetImplemented
}

// Eopkg specific functions

func (e *EopkgManager) eopkgExecRoot(args []string) error {
	endArgs := []string{
		"-D", e.root,
	}
	args = append(args, endArgs...)
	return ExecStdoutArgs("eopkg", args)
}

// Add a repository to the target
func (e *EopkgManager) addRepos(ops []spec.Operation) error {
	for _, repo := range ops {
		r := repo.(*spec.OpRepo)
		if err := e.eopkgExecRoot([]string{"add-repo", r.RepoName, r.RepoURI}); err != nil {
			return err
		}
	}
	return nil
}

// Install the named components
func (e *EopkgManager) installComponents(ops []spec.Operation, ignoreSafety bool) error {
	var componentNames []string
	for _, comp := range ops {
		c := comp.(*spec.OpGroup)
		componentNames = append(componentNames, c.GroupName)
	}
	cmd := []string{"install", "-y", "--ignore-comar", "-c"}
	cmd = append(cmd, componentNames...)
	if ignoreSafety {
		cmd = append(cmd, "--ignore-safety")
	}
	return e.eopkgExecRoot(cmd)
}

// Install the named packages
func (e *EopkgManager) installPackages(ops []spec.Operation, ignoreSafety bool) error {
	var pkgNames []string
	for _, p := range ops {
		pk := p.(*spec.OpPackage)
		pkgNames = append(pkgNames, pk.Name)
	}
	cmd := []string{"install", "-y", "--ignore-comar"}
	cmd = append(cmd, pkgNames...)
	if ignoreSafety {
		cmd = append(cmd, "--ignore-safety")
	}
	return e.eopkgExecRoot(cmd)
}
