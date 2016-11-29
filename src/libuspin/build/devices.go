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
	"libosdev/commands"
	"path/filepath"
)

// DeviceNode represents a /dev/ node to be created in chroots
type DeviceNode struct {
	Mode  string // Mode to create the device node in
	Major uint32 // Major ID
	Minor uint32 // Minor ID
	Path  string // Path within a chroot (no / prefix)
}

var (
	// DevNodeRandom is /dev/random
	DevNodeRandom *DeviceNode

	// DevNodeURandom is /dev/urandom
	DevNodeURandom *DeviceNode
)

func init() {
	DevNodeURandom = &DeviceNode{Mode: "00666", Major: 1, Minor: 9, Path: "dev/urandom"}
	DevNodeRandom = &DeviceNode{Mode: "00666", Major: 1, Minor: 8, Path: "dev/random"}
}

// CreateDeviceNode will create the essential nodes in a chroot path
func CreateDeviceNode(root string, node *DeviceNode) error {
	fpath := filepath.Join(root, node.Path)
	cmd := []string{"-m", node.Mode, fpath, "c", fmt.Sprintf("%d", node.Major), fmt.Sprintf("%d", node.Minor)}

	return commands.ExecStdoutArgs("mknod", cmd)
}
