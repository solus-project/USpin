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

// InstallPackages will install all required packages into the rootfs
func (s *SolSpin) InstallPackages() error {
	s.logPackage.Info("Applying operations")

	// First thing first, ensure that it always cleans up within this context,
	// so that we know it's done before we return to building the image
	defer func() {
		if err := s.packager.Cleanup(); err != nil {
			s.logPackage.Error(err)
		}
	}()

	// Attempt to init root now
	s.logPackage.Info("Initialising root with package manager")
	if err := s.packager.InitRoot(s.builder.GetRootDir()); err != nil {
		return err
	}

	for _, opset := range s.spec.Stack.Blocks {
		if err := s.packager.ApplyOperations(opset.Ops); err != nil {
			return err
		}
	}

	s.logPackage.Info("Finalizing package operations")
	if err := s.packager.FinalizeRoot(); err != nil {
		return err
	}
	return nil
}
