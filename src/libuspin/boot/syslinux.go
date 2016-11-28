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

// SyslinuxLoader wraps isolinux/syslinux into a single set of management
// routines
type SyslinuxLoader struct {
}

// Init will attempt to initialise this loader if all host requirements are
// actually met.
func (s *SyslinuxLoader) Init() error {
	return ErrNotYetImplemented
}

// NewSyslinuxLoader will return a newly created SyslinuxLoader instance
func NewSyslinuxLoader() *SyslinuxLoader {
	return &SyslinuxLoader{}
}
