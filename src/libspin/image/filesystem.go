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
	"github.com/Sirupsen/logrus"
)

var filesystemCommands map[string]string

func init() {
	// Initialise the filesystemCommands
	filesystemCommands = make(map[string]string)
	filesystemCommands["ext4"] = "mkfs -t ext4 -F %s"
}

// FormatAs will format the given path with the filesystem specified.
// Note: You should only use this with image paths, it's dangerous!
func FormatAs(filename string, filesystem string) error {
	command, ok := filesystemCommands[filesystem]
	if !ok {
		return fmt.Errorf("Cannot format with unknown filesystem '%v'", filesystem)
	}
	log.WithFields(logrus.Fields{
		"filename":   filename,
		"filesystem": filesystem,
	}).Info("Formatting filesystem")
	return ExecStdout(fmt.Sprintf(command, filename))
}
