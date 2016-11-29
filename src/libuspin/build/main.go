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

// Package build provides routines for the core elements of image building
package build

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"libuspin/commands"
	"libuspin/config"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// A Kernel is exactly what it looks like. :p
type Kernel struct {
	Version  string
	Path     string
	BaseName string
}

var (
	// ErrNoKernelFound is returned when a builder cannot find a kernel in the given root
	ErrNoKernelFound = errors.New("Could not find a valid kernel")
)

// CreateSparseFile will create a new sparse file with the given filename and
// size in nMegabytes/
//
// This is highly dependent on the underlying filesystem at the directory
// where the file is to be created, making use of the syscall ftruncate.
func CreateSparseFile(filename string, nMegabytes int) error {
	log.WithFields(log.Fields{
		"filename": filename,
		"size":     nMegabytes}).Info("Creating sparse file")
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 00644)
	if err != nil {
		return err
	}
	defer f.Close()
	// NOTE: New megabytes, not old megabytes (1000, not 1024)
	sz := int64(nMegabytes * 1000 * 1000)
	if err = syscall.Ftruncate(int(f.Fd()), sz); err != nil {
		return err
	}
	return nil
}

// GetSquashfsArgs returns the compression arg set for a given compression type
func GetSquashfsArgs(compressionType config.CompressionType) ([]string, error) {
	switch compressionType {
	case config.CompressionGzip:
		return []string{"-comp", "gzip"}, nil
	case config.CompressionXZ:
		return []string{"-comp", "xz"}, nil
	default:
		return nil, fmt.Errorf("Unknown compression type: %v", compressionType)
	}
}

// CreateSquashfs will create a new squashfs filesystem image at the given outputFile path,
// containing the tree found at path, using compressionType (gzip or xz).
func CreateSquashfs(path, outputFile string, compressionType config.CompressionType) error {
	command := []string{
		path,
		outputFile,
	}
	dirName := ""
	if fp, err := filepath.Abs(path); err == nil {
		dirName = filepath.Dir(fp)
	} else {
		return err
	}

	// May have to set -keep-as-directory
	if st, err := os.Stat(path); err == nil {
		if st.Mode().IsDir() {
			command = append(command, "-keep-as-directory")
		}
	} else {
		return err
	}
	// TODO: Add -processors nCPU (4)
	if execArgs, err := GetSquashfsArgs(compressionType); err == nil {
		command = append(command, execArgs...)
	} else {
		return err
	}
	return commands.ExecStdoutArgsDir(dirName, "mksquashfs", command)
}

// GetKernelFromRoot will attempt to "learn" about the kernel from the rootfs
// and return a populated kernel struct.
func GetKernelFromRoot(root string) (*Kernel, error) {
	// Add more as time goes by.
	possiblePaths := []string{
		filepath.Join(root, "vmlinuz"),
		filepath.Join(root, "boot", "vmlinuz"),
	}

	for _, p := range possiblePaths {
		// as an example, /vmlinuz -> boot/kernel-4.8.10
		st, err := os.Stat(p)
		if err != nil || st == nil {
			continue
		}

		// TODO: Add more scan thingers.
		baseNom := filepath.Base(p)
		splits := strings.Split(p, "-")
		if len(splits) < 2 {
			log.WithFields(log.Fields{
				"kernel": baseNom,
			}).Warning("Don't know how to handle kernel version")
			continue
		}
		version := strings.Join(splits[1:], "-")
		log.WithFields(log.Fields{
			"kernel":  baseNom,
			"version": version,
		}).Info("Discovered usable kernel")
		return &Kernel{
			Version:  version,
			BaseName: baseNom,
			Path:     p,
		}, nil
	}

	return nil, ErrNoKernelFound
}
