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

package config

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

// ImageConfiguration is the configuration for an image build
type ImageConfiguration struct {
}

// New will return a new ImageConfiguration for the given path and attempt to
// parse it. This function will return a nil ImageConfiguration if parsing
// fails.
func New(cpath string) (*ImageConfiguration, error) {
	iconf := &ImageConfiguration{}
	var data []byte
	var err error
	var fi *os.File

	fi, err = os.Open(cpath)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	// Read the configuration file in
	if data, err = ioutil.ReadAll(fi); err != nil {
		return nil, err
	}

	// Attempt to populate config from the toml spin file
	if _, err = toml.Decode(string(data), iconf); err != nil {
		return nil, err
	}
	return iconf, nil
}
