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

// Package libimage provides routines for the core elements of image building
package libimage

import (
	"os"
	"syscall"
)

// CreateSparseFile will create a new sparse file with the given filename and
// size in nMegabytes/
//
// This is highly dependent on the underlying filesystem at the directory
// where the file is to be created, making use of the syscall ftruncate.
func CreateSparseFile(filename string, nMegabytes int) error {
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
