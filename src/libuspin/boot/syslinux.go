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

import (
	"fmt"
	"libuspin/config"
	"libuspin/disk"
	"os"
	"path/filepath"
	"text/template"
)

// IsolinuxTemplate is used to populate fields in the isolinux.cfg
type IsolinuxTemplate struct {
	Kernel      *Kernel
	Label       string // CDLABEL
	Title       string // Needs to come from config!
	StartString string
}

var (
	// DefaultIsolinuxTemplate is the built-in template for isolinux.cfg
	DefaultIsolinuxTemplate = `
ui vesamenu.c32
timeout 50
default live

MENU RESOLUTION 1024 768
menu title {{.Title}}

menu color screen       37;40      #80ffffff #00000000 std
MENU COLOR border       30;44   #40ffffff #a0000000 std
MENU COLOR title        1;36;44 #ffffffff #a0000000 std
MENU COLOR sel          7;37;40 #e0ffffff #20ffffff all
MENU COLOR unsel        37;44   #50ffffff #a0000000 std
MENU COLOR help         37;40   #c0ffffff #a0000000 std
MENU COLOR timeout_msg  37;40   #80ffffff #00000000 std
MENU COLOR timeout      1;37;40 #c0ffffff #00000000 std
MENU COLOR msg07        37;40   #90ffffff #a0000000 std
MENU COLOR tabmsg       31;40   #ffDEDEDE #00000000 std
MENU HIDDEN
MENU HIDDENROW 7
MENU WIDTH 78
MENU MARGIN 15
MENU ROWS 4
MENU VSHIFT 17
MENU HSHIFT 25
MENU TABMSGROW 11

label live
  menu label {{.StartString}}
  kernel /{{.Kernel.TargetPath}}
  append initrd=/{{.Kernel.TargetInitrd}} root=live:CDLABEL={{.Label}} ro rd.luks=0 rd.md=0 quiet splash --
menu default
label local
  menu label Boot from local drive
  localboot 0x80
`
)

var (
	// SyslinuxPaths contains paths known to be used by the majority of Linux
	// distributions, so that we can search for required files prior to actually
	// trying to spin an ISO.
	SyslinuxPaths = []string{
		"/usr/lib64/syslinux",
		"/usr/lib/syslinux",
		"/usr/share/syslinux",
	}

	// SyslinuxAssets is the core set of assets required by all syslinux usages
	SyslinuxAssets = []string{
		"libutil.c32",
		"libcom32.c32",
		"ldlinux.c32",
	}

	// SyslinuxAssetsISO are the assets required explicitly for ISOs, i.e. menu bits
	SyslinuxAssetsISO = []string{
		"vesamenu.c32",
		"isolinux.bin",
		"vesa.c32",
		"isohdpfx.bin",
	}
)

// SyslinuxLoader wraps isolinux/syslinux into a single set of management
// routines
type SyslinuxLoader struct {
	// A basename to full path mapping of the asset paths (i.e. vesamenu.c32 -> /usr/share/blah)
	cachedAssets map[string]string

	// Store the configuration for particulars we need to implement
	config *config.ImageConfiguration

	isolinuxTemplate *template.Template
}

// LocateAsset will attempt to find the given asset and then cache it
func (s *SyslinuxLoader) LocateAsset(name string) error {
	if _, ok := s.cachedAssets[name]; ok {
		return nil
	}
	for _, path := range SyslinuxPaths {
		fpath := filepath.Join(path, name)
		if _, err := os.Stat(fpath); err != nil {
			continue
		}
		s.cachedAssets[name] = fpath
		return nil
	}
	return fmt.Errorf("Cannot find required syslinux asset: %v", name)
}

// Init will attempt to initialise this loader if all host requirements are
// actually met.
func (s *SyslinuxLoader) Init(c *config.ImageConfiguration) error {
	for _, item := range SyslinuxAssets {
		if err := s.LocateAsset(item); err != nil {
			return err
		}
	}
	// TODO: Only call this when producing an ISO!
	for _, item := range SyslinuxAssetsISO {
		if err := s.LocateAsset(item); err != nil {
			return err
		}
	}
	// Parse our default templates
	// TODO: Make it possible to load from a file!
	tmpl, err := template.New("isolinux").Parse(DefaultIsolinuxTemplate)
	if err != nil {
		return err
	}
	s.isolinuxTemplate = tmpl
	s.config = c
	return nil
}

// GetCapabilities will return isolinux support only for syslinux right now
func (s *SyslinuxLoader) GetCapabilities() Capability {
	return CapInstallISO | CapInstallLegacy
}

// NewSyslinuxLoader will return a newly created SyslinuxLoader instance
func NewSyslinuxLoader() *SyslinuxLoader {
	s := &SyslinuxLoader{
		cachedAssets: make(map[string]string),
	}
	return s
}

// Install will do the real work of installing syslinux bootloader
func (s *SyslinuxLoader) Install(op Capability, c ConfigurationSource) error {
	// Currently we're only ever invoked as Legacy|ISO
	bootdirTarget := c.JoinDeployPath("isolinux")

	// First off actually try to install the boot directory
	if err := os.MkdirAll(bootdirTarget, 00755); err != nil {
		return err
	}

	// Collect required assets
	var reqAssets []string
	for _, key := range SyslinuxAssets {
		reqAssets = append(reqAssets, key)
	}
	for _, key := range SyslinuxAssetsISO {
		reqAssets = append(reqAssets, key)
	}

	// Install the ISO assets
	for _, asset := range reqAssets {
		target := c.JoinDeployPath("isolinux", asset)
		if err := disk.CopyFile(s.cachedAssets[asset], target); err != nil {
			return err
		}
	}

	brand := s.config.Branding.Title
	str := s.config.Branding.StartString

	// Write our template data
	tmplData := IsolinuxTemplate{
		Kernel:      c.GetKernel(),
		Label:       "DummyISO",
		Title:       brand,
		StartString: str,
	}

	cfg := c.JoinDeployPath("isolinux", "isolinux.cfg")
	out, err := os.Create(cfg)
	if err != nil {
		return err
	}
	defer out.Close()
	return s.isolinuxTemplate.Execute(out, tmplData)
}

// GetSpecialFile will return the special paths for isolinux
func (s *SyslinuxLoader) GetSpecialFile(t FileType) string {
	// Currently we're only ever invoked as Legacy|ISO
	switch t {
	case FileTypeBootElToritoBinary:
		return filepath.Join("isolinux", "isolinux.bin")
	case FileTypeBootElToritoCatalog:
		return filepath.Join("isolinux", "boot.cat")
	case FileTypeBootMBR:
		return filepath.Join("isolinux", "isohdpfx.bin")
	default:
		return ""
	}
}
