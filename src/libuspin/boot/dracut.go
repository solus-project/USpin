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
	"github.com/solus-project/libosdev/commands"
	"strings"
)

var (
	// DracutLiveOSModules are modules to enable for Live OS Usage
	// TODO: Make systemd dependent on the presence of systemd in the sysroot
	DracutLiveOSModules = []string{"dmsquash-live", "systemd", "pollcdrom"}

	// DracutLiveOSDrivers are drivers that should be shipped for LiveOS functionality to work
	// TODO: Investigate now-dead stuff and curate this list
	DracutLiveOSDrivers = []string{
		"squashfs",
		"ext2",
		"vfat",
		"msdos",
		"sr_mod",
		"sd_mod",
		"ehci_hcd",
		"uhci_hcd",
		"xhci_hcd",
		"xhci_pci",
		"ohci_hcd",
		"usb_storage",
		"usbhid",
		"dm_mod",
		"ata_generic",
		"libata",
	}
)

// Dracut provides wrapping around dracut generation in chroots
type Dracut struct {
	// Additional options to pass to dracut when generating
	Options []string

	// Extra modules to enable
	Modules []string

	// Extra drivers to enable
	Drivers []string

	// The filename to use within the root (should include / prefix)
	OutputFilename string

	// Compression method, i.e. --lz4, --gzip, etc.
	CompressionMethod string

	k *Kernel
}

// NewDracut returns a new Dracut object for generation
func NewDracut(k *Kernel) *Dracut {
	return &Dracut{
		OutputFilename:    fmt.Sprintf("/boot/initramfs-%v.img", k.Version),
		CompressionMethod: "--lz4",
		k:                 k,
	}
}

// Exec will build the command run dracut within the chroot
func (d *Dracut) Exec(path string) error {
	cmd := fmt.Sprintf("dracut --no-hostonly-cmdline -N --kver \"%v\"", d.k.Version)

	if d.CompressionMethod != "" {
		cmd += " " + d.CompressionMethod
	}

	if len(d.Modules) > 0 {
		cmd += fmt.Sprintf(" --add \"%v\"", strings.Join(d.Modules, " "))
	}
	if len(d.Drivers) > 0 {
		cmd += fmt.Sprintf(" --add-drivers \"%v\"", strings.Join(d.Drivers, " "))
	}

	if !strings.HasPrefix(d.OutputFilename, "/") {
		return fmt.Errorf("Invalid dracut name: %v", d.OutputFilename)
	}

	cmd += fmt.Sprintf(" \"%v\"", d.OutputFilename)

	return commands.ChrootExec(path, cmd)
}
