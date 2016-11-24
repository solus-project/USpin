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

package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"libspin"
	"libspin/image"
	"libspin/pkg"
	"os"
)

var log *logrus.Logger

func init() {
	form := &logrus.TextFormatter{}
	form.FullTimestamp = true
	form.TimestampFormat = "15:04:05"
	log = logrus.New()
	log.Out = os.Stderr
	log.Formatter = form
}

func printUsage(exitCode int) {
	var fd *os.File
	if exitCode == 0 {
		fd = os.Stdout
	} else {
		fd = os.Stderr
	}

	fmt.Fprintf(fd, "%s [image.spin]\n", os.Args[0])
	os.Exit(exitCode)
}

func doBuild(spinfile string) {
	var spec *libspin.ImageSpec
	var err error
	var builder image.Builder
	var pman pkg.Manager

	// Check the user has root privs
	if os.Geteuid() != 0 {
		log.WithFields(logrus.Fields{"type": "privilege", "euid": os.Geteuid()}).Error("solspin requires root privileges")
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{"filename": spinfile}).Info("Loading .spin file")

	if spec, err = libspin.NewImageSpec(spinfile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// For nicer logging
	imgType := spec.Config.Image.Type
	// TODO: Don't hardcode package manager??
	packageManager := "eopkg"
	logImg := log.WithFields(logrus.Fields{"imageType": imgType})
	logPkg := log.WithFields(logrus.Fields{"packageManager": packageManager})

	// Get the builder instance
	builder, err = image.NewBuilder(spec.Config.Image.Type)
	if err != nil {
		logImg.Fatal(err)
	}
	defer builder.Cleanup()

	logImg.Info("Initialising builder")
	if err = builder.Init(spec); err != nil {
		logImg.Error(err)
		return
	}

	// Try to init the package manager
	pman, err = pkg.NewManager("eopkg")
	if err != nil {
		logPkg.Error(err)
		return
	}

	logPkg.Info("Initialising package manager")
	if err = pman.Init(spec.Config); err != nil {
		logPkg.Error(err)
		return
	}

	logImg.Info("Preparing workspace")
	if err = builder.PrepareWorkspace(); err != nil {
		logImg.Error(err)
		return
	}

	logImg.Info("Creating storage")
	if err = builder.CreateStorage(); err != nil {
		logImg.Error(err)
		return
	}

	logImg.Info("Mounting storage")
	if err = builder.MountStorage(); err != nil {
		logImg.Error(err)
		return
	}

	// Attempt to init root now
	logPkg.Info("Initialising root with package manager")
	if err = pman.InitRoot(builder.GetRootDir()); err != nil {
		logPkg.Error(err)
		return
	}

	logPkg.Info("Applying operations")
	for _, opset := range spec.Stack.Blocks {
		if err := pman.ApplyOperations(opset.Ops); err != nil {
			logPkg.Error(err)
			return
		}
	}

	logPkg.Info("Finalizing package operations")
	if err = pman.FinalizeRoot(); err != nil {
		logPkg.Error(err)
		return
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage(1)
	}
	doBuild(os.Args[1])
	log.Error("Not yet implemented")
}
