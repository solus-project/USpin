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

// Package config provides the TOML based configuration format used for the
// main .spin files
package config

import (
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"strings"
)

// ImageType is the type of image that will be created
type ImageType string

// A LoaderType is a pseudo enum type for the bootloader to restrict to
// supported implementations
type LoaderType string

const (
	// ImageTypeLiveOS is an ISO type image that may also be USB compatible
	ImageTypeLiveOS ImageType = "liveos"
)

const (
	// LoaderTypeSyslinux refers to syslinux + isolinux
	LoaderTypeSyslinux LoaderType = "syslinux"
)

// SectionImage describes the [image] portion of a spin file
type SectionImage struct {
	Packages string    `toml:"packages"` // Path to the packages file
	Type     ImageType `toml:"type"`     // Type of image to construct
}

// SectionBranding describes the image branding rules
type SectionBranding struct {
	Title string `toml:"title"` // Title of the OS to use in bootloaders
}

// ImageConfiguration is the configuration for an image build
type ImageConfiguration struct {
	Image    SectionImage    `toml:"image"`
	Branding SectionBranding `toml:"branding"`
	LiveOS   SectionLiveOS   `toml:"liveos"`
	Isolinux SectionIsolinux `toml:"isolinux"`
}

// New will return a new ImageConfiguration for the given path and attempt to
// parse it. This function will return a nil ImageConfiguration if parsing
// fails.
func New(cpath string) (*ImageConfiguration, error) {
	iconf := &ImageConfiguration{
		LiveOS: SectionLiveOS{
			RootfsFormat: "ext4",
			RootfsSize:   4000,
			BootDir:      "boot",
			// Default to isolinux
			Bootloaders: []LoaderType{
				LoaderTypeSyslinux,
			},
		},
	}
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

	// Ensure errors is non empty!
	iconf.Image.Packages = strings.TrimSpace(iconf.Image.Packages)
	if iconf.Image.Packages == "" {
		return nil, errors.New("image.packages cannot be empty")
	}

	// Validate the type
	// TODO: Add more image types!
	switch iconf.Image.Type {
	case ImageTypeLiveOS:
		if err := ValidateSectionLiveOS(&iconf.LiveOS); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Unknown image type: %v", iconf.Image.Type)
	}

	return iconf, nil
}
