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
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var filesystemCommands map[string]string

var log *logrus.Logger

func init() {
	// Initialise the filesystemCommands
	filesystemCommands = make(map[string]string)
	filesystemCommands["ext4"] = "mkfs -t ext4 -F %s"

	// Create the logger
	form := &logrus.TextFormatter{}
	form.FullTimestamp = true
	form.TimestampFormat = "15:04:05.00"
	log = logrus.New()
	log.Out = os.Stderr
	log.Formatter = form
}

// ExecStdout is a convenience function to execute a command to the stdout
// and return the error, if any
func ExecStdout(command string) error {
	splits := strings.Fields(command)
	var c *exec.Cmd
	cmdName := splits[0]
	var err error
	// Search the path if necessary
	if !strings.Contains(cmdName, "/") {
		cmdName, err = exec.LookPath(cmdName)
		if err != nil {
			return err
		}
	}
	// Ensure we pass arguments
	if len(splits) == 1 {
		c = exec.Command(cmdName)
	} else {
		c = exec.Command(cmdName, splits[1:]...)
	}
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

// CreateSparseFile will create a new sparse file with the given filename and
// size in nMegabytes/
//
// This is highly dependent on the underlying filesystem at the directory
// where the file is to be created, making use of the syscall ftruncate.
func CreateSparseFile(filename string, nMegabytes int) error {
	log.WithFields(logrus.Fields{
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
