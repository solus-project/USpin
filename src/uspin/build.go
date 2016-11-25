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

// Build will attempt to build the image, and return an error if this fails
func (s *SolSpin) Build() error {
	// Initialise our builder before we go anywhere
	if err := s.builder.Init(s.spec); err != nil {
		s.logImage.Error(err)
		return err
	}
	// Make sure that the package manager requirements are met
	if err := s.packager.Init(s.spec.Config); err != nil {
		s.logPackage.Error(err)
		return err
	}

	// Always perform cleanup duty.
	defer s.builder.Cleanup()

	// Start building the base parts of the image
	if err := s.StartImageBuild(); err != nil {
		s.logImage.Error(err)
		return err
	}

	// Hand over to the package manager
	if err := s.InstallPackages(); err != nil {
		s.logPackage.Error(err)
		return err
	}

	// And now finish the image build
	if err := s.FinishImageBuild(); err != nil {
		s.logImage.Error(err)
		return err
	}

	// TODO: Finish the image
	return nil
}
