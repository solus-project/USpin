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

package build

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"libuspin"
	"libuspin/boot"
	"libuspin/commands"
	"libuspin/disk"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	// ErrNotYetImplemented is just used until we actually implement some code....
	ErrNotYetImplemented = errors.New("Not yet implemented")

	requiredBinaries []string
)

const (
	// DefaultImageSize is the size of the rootfs we try to create (4GB)
	DefaultImageSize = 4000
)

func init() {
	requiredBinaries = []string{
		"isohybrid",
		"mksquashfs",
		"xorriso",
	}
}

// A LiveOSBuilder is responsible for building ISO format images that are USB
// compatible. It is the "LiveCD" type of Builder
type LiveOSBuilder struct {
	img            *libuspin.ImageSpec
	rootfsImg      string
	rootfsDir      string
	rootfsFormat   string
	rootfsSize     int
	deployDir      string
	liveosDir      string
	liveStagingDir string
	workspace      string

	cdlabel string // What to name the ISO

	// For storing bootloader bits
	loaders []boot.Loader

	// The kernel to be used for booting
	kernel *boot.Kernel
}

// NewLiveOSBuilder should only be used by builder.go
func NewLiveOSBuilder() *LiveOSBuilder {
	// TODO: Unhardcode!
	return &LiveOSBuilder{
		cdlabel: "DummyISO",
	}
}

// Init will initialise a LiveOSBuilder from the given spec
func (l *LiveOSBuilder) Init(img *libuspin.ImageSpec) error {
	l.img = img

	// Ensure all required binaries are available before we go doing anything.
	for _, bin := range requiredBinaries {
		if _, err := exec.LookPath(bin); err != nil {
			return err
		}
	}

	// rootfs.img particulars
	l.rootfsFormat = l.img.Config.LiveOS.RootfsFormat
	l.rootfsSize = l.img.Config.LiveOS.RootfsSize

	// Init the bootloaders
	if loaders, err := boot.InitLoaders(l.img.Config, l.img.Config.LiveOS.Bootloaders); err == nil {
		l.loaders = loaders
	} else {
		return err
	}

	// Primary bootloader is one with Legacy&ISO capability
	if !boot.HaveLoaderWithMask(l.loaders, boot.CapInstallISO|boot.CapInstallLegacy) {
		return errors.New("No usable bootloader found. Need ISO|Legacy")
	}

	return nil
}

// JoinPath is a helper to join paths onto our root workspace directory
func (l *LiveOSBuilder) JoinPath(paths ...string) string {
	return filepath.Join(l.workspace, filepath.Join(paths...))
}

// PrepareWorkspace sets up the required directories for the LiveOSBuilder
func (l *LiveOSBuilder) PrepareWorkspace() error {
	var err error
	if l.workspace, err = filepath.Abs("./workspace"); err != nil {
		return err
	}

	// Purge existing workspace always
	if st, err := os.Stat(l.workspace); err == nil {
		if st != nil && st.Mode().IsDir() {
			if err = os.RemoveAll(l.workspace); err != nil {
				return err
			}
		}
	}

	// Initialise our base variables
	l.rootfsDir = l.JoinPath("rootfs")
	l.deployDir = l.JoinPath("deploy")
	// Inside the ISO target
	l.liveosDir = l.JoinPath("deploy", "LiveOS")
	// Inside the workspace only
	l.liveStagingDir = l.JoinPath("LiveOS")
	l.rootfsImg = l.JoinPath("LiveOS", "rootfs.img")

	// As and when we add new directories, populate them here
	requiredDirs := []string{
		l.workspace,
		l.rootfsDir,
		l.deployDir,
		l.liveosDir,
		l.liveStagingDir,
	}

	// Create all required directories
	for _, dir := range requiredDirs {
		if err = os.MkdirAll(dir, 00755); err != nil {
			return err
		}
	}

	return nil
}

// CreateStorage will create the rootfs.img in which we will contain the
// Live OS
func (l *LiveOSBuilder) CreateStorage() error {
	if err := CreateSparseFile(l.rootfsImg, l.rootfsSize); err != nil {
		return err
	}
	if err := FormatAs(l.rootfsImg, l.rootfsFormat); err != nil {
		return err
	}
	return nil
}

// Cleanup currently does nothing within this builder
func (l *LiveOSBuilder) Cleanup() {
	log.Info("Cleaning up")
	disk.GetMountManager().UnmountAll()
}

// MountStorage will mount the rootfs.img so that the package manager can
// take over
func (l *LiveOSBuilder) MountStorage() error {
	return disk.GetMountManager().Mount(l.rootfsImg, l.rootfsDir, l.rootfsFormat, "loop")
}

// UnmountStorage will unmount the rootfs.img from earlier
// This is the last point in which the storage is used, so we check the filesystem
// is OK here.
func (l *LiveOSBuilder) UnmountStorage() error {
	if err := disk.GetMountManager().Unmount(l.rootfsDir); err != nil {
		return err
	}
	return CheckFS(l.rootfsImg, l.rootfsFormat)
}

// GetRootDir returns the path to the mounted rootfs.img
func (l *LiveOSBuilder) GetRootDir() string {
	return l.rootfsDir
}

// The very last call in the chain, we seal the deal by spinning the ISO
func (l *LiveOSBuilder) spinISO() error {
	uefi := false
	// Get absolute path for "./${name}"
	outputFilename := l.img.Config.LiveOS.FileName
	if o, err := filepath.Abs(outputFilename); err == nil {
		outputFilename = o
	} else {
		return err
	}
	volumeID := l.cdlabel
	command := []string{
		"-no_rc", // Forbid reading startup files which may skew ISO generation
		"-as",
		"mkisofs",
		"-iso-level",
		"3",
		"-full-iso9660-filenames",
		"-volid",
		volumeID,
		"-appid",
		volumeID,
	}

	caps := boot.CapInstallISO | boot.CapInstallLegacy
	bloader := boot.GetLoaderWithMask(l.loaders, caps)

	bootbinFile := bloader.GetSpecialFile(boot.FileTypeBootElToritoBinary)
	bootcatFile := bloader.GetSpecialFile(boot.FileTypeBootElToritoCatalog)
	mbrFile := bloader.GetSpecialFile(boot.FileTypeBootMBR)

	// This is where we'd install syslinux or other loader..
	// Note we'll need to do investigation for GRUB to determine precisely how to
	// get the cat and bin files
	if bootbinFile != "" && bootcatFile != "" {
		command = append(command, []string{
			"-eltorito-boot",
			bootbinFile,
			"-eltorito-catalog",
			bootcatFile,
			"-no-emul-boot",
			"-boot-load-size",
			"4",
			"-boot-info-table",
		}...)
	}
	// Enable USB booting per the bootloader
	if mbrFile != "" {
		command = append(command, []string{
			"-isohybrid-mbr",
			mbrFile,
		}...)
	}
	// Still unused for now
	if uefi {
		command = append(command, []string{
			"-eltorito-alt-boot",
			"-e",
			"efi.img", // TODO: Use appropriate origin
			"-no-emul-boot",
			"-isohybrid-gpt-basdat",
		}...)
	}
	// Set the output filename and directory
	command = append(command, []string{
		"-output",
		outputFilename,
		".", // Create from current directory
	}...)
	return commands.ExecStdoutArgsDir(l.deployDir, "xorriso", command)
}

// Install the bootloader for the given image
func (l *LiveOSBuilder) installBootloader() error {
	// For now we're only interested in legacy install, so deal with that
	caps := boot.CapInstallISO | boot.CapInstallLegacy
	bloader := boot.GetLoaderWithMask(l.loaders, caps)

	return bloader.Install(caps, l)
}

// CollectAssets will collect the kernel and create a new initramfs to be used
// during the boot process
func (l *LiveOSBuilder) CollectAssets() error {
	kernel, err := boot.GetKernelFromRoot(l.rootfsDir)
	if err != nil {
		return err
	}
	l.kernel = kernel

	// Create the boot/ directory
	bootbase := l.img.Config.LiveOS.BootDir
	bootdir := l.JoinDeployPath(bootbase)
	if err := os.MkdirAll(bootdir, 00755); err != nil {
		return err
	}

	// Copy the kernel, standard "kernel" name
	ktgt := filepath.Join(bootdir, "kernel")
	if err := disk.CopyFile(kernel.Path, ktgt); err != nil {
		return err
	}

	// Required by the bootloaders
	l.kernel.TargetPath = filepath.Join(bootbase, "kernel")
	l.kernel.TargetInitrd = filepath.Join(bootbase, "initrd.img")

	// Attempt to build dracut image
	drac := boot.NewDracut(l.kernel)
	drac.Modules = boot.DracutLiveOSModules
	drac.Drivers = boot.DracutLiveOSDrivers
	drac.OutputFilename = "/live.img"

	if err := drac.Exec(l.rootfsDir); err != nil {
		return err
	}

	// Copy the new live.img asset across
	dracSource := filepath.Join(l.rootfsDir, "live.img")
	dracTarget := filepath.Join(bootdir, "initrd.img")
	if err := disk.CopyFile(dracSource, dracTarget); err != nil {
		return err
	}

	// Nuke live.img from the filesystem
	if err := os.Remove(dracSource); err != nil {
		return err
	}

	return nil
}

// FinalizeImage will go ahead and finish up the ISO construction
func (l *LiveOSBuilder) FinalizeImage() error {
	// First up, create the squashfs
	squash := filepath.Join(l.liveosDir, "squashfs.img")
	if err := CreateSquashfs(l.liveStagingDir, squash, l.img.Config.LiveOS.Compression); err != nil {
		return err
	}

	// Attempt installation of bootloader
	if err := l.installBootloader(); err != nil {
		return err
	}

	// TODO: Install bootloader, copy asset files, put kernel in place, etc.
	return l.spinISO()
}

//
// The following are all ConfigurationSource methods
//

// GetBootDevice always returns nil for LiveOS
func (l *LiveOSBuilder) GetBootDevice() string {
	return ""
}

// GetRootDevice will actually return the cdlabel for ISO mode bootloaders
func (l *LiveOSBuilder) GetRootDevice() string {
	return l.cdlabel
}

// JoinDeployPath will return a path within the LiveOS workspace
func (l *LiveOSBuilder) JoinDeployPath(paths ...string) string {
	return filepath.Join(l.deployDir, filepath.Join(paths...))
}

// JoinRootPath will return a path within the LiveOS rootfs image
func (l *LiveOSBuilder) JoinRootPath(paths ...string) string {
	return filepath.Join(l.rootfsDir, filepath.Join(paths...))
}

// GetKernel returns our stored kernel object
func (l *LiveOSBuilder) GetKernel() *boot.Kernel {
	return l.kernel
}
