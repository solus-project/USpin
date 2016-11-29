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

package boot

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
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
		if kpath, err := filepath.EvalSymlinks(p); err == nil {
			p = kpath
		} else {
			return nil, err
		}

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
