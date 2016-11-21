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

package libspin

// ImageSpecParser does the heavy lifting of parsing a .spin file to pull all
// relevant stack operations from it.
type ImageSpecParser struct {
}

// NewParser will return a new parser for the image specification file
func NewParser(path string) *ImageSpecParser {
	return &ImageSpecParser{}
}

// Parse will attempt to parse the given image speicifcation file at the given
// path, and will return an error if this fails.
func (i *ImageSpecParser) Parse(path string) error {
	return errors.New("Not yet implemented!")
}
