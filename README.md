USpin
------

[![Report](https://goreportcard.com/badge/github.com/solus-project/USpin)](https://goreportcard.com/report/github.com/solus-project/USpin) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Universal image creation utilities. Used to produce bootable Linux medium in an agnostic fashion. This is still a **Work In Progress**.

USpin is a [Solus project](https://getsol.us/).

![logo](https://build.getsol.us/logo.png)


[![You Spin Me Right Round](https://img.youtube.com/vi/PGNiXGX2nLU/0.jpg)](https://www.youtube.com/watch?v=PGNiXGX2nLU)

**Note**
Solus, USpin, and the Solus Project Team, are not endorsed by any projects listed here in terms of planned support. We are simply creating a tool to simplify
management and creation of various images all from one place. As developers may experiment with, be involved in, or test other projects and distros for
many reasons (including validation of projects we develop such as Budgie) - it is useful even to us to be able to produce medium for other distros using
a single standardized tool.

Obviously the core benefit to us is in producing our own medium, however others are welcome to join in and contribute too, because there are a great many
tools for creating distro images, each one more specialised than the last. It also serves somewhat as a research project, in seeing "how others do" to find
ways in which Solus can improve itself too. The more specialised the support for Solus in the tool in contrast to other distros will highlight to us exactly
what we can and should improve.


Rationale
-------
Intended to succeed the existing `solus-image-creator.py` script with something a bit more robust that can construct multiple image types.

Currently the existing image creator can only construct a simple `ISO9660` image, however Solus also makes use of chrootable base images for the `overlayfs` system employed in `evobuild`.

Planned support
---------------

Those with a fire symbol (🔥) are currently considered important to Solus projects and are the active priority. After stabilisation we can work on the support for other medium, and at that point will happily take contributions (post v1)

**Medium**

 - `LiveOS` (dracut distros) 🔥
 - `raw` filesystem type (partitions in disk image) 🔥
 - `flat` image support (no partitions, i.e. an `ext4` loopback image) 🔥
 - `casper` (Ubuntu)
 - `debian-live` (Vanilla Debian images)

**Boot**:
 - `isolinux`/`syslinux` 🔥
 - `systemd`-class bootloaders 🔥
 - `grub` "2"

**Package Management**:

- `eopkg` (done) 🔥
- `sol` (for validation in Solus)
- `yum`
- `dnf`
- `swupd`
- `.deb` (`dpkg`/`apt-get`/`apt`) (via `debootstrap` maybe?)

TODO
----

 - [x] Add parser for the Solus image specification format
 - [x] Port the `Stack` implementation from old image creator
 - [x] Add config format for the main image configuration
 - [x] Add utilities for image format & creation (`dd`/`fallocate`, etc)
 - [x] Implement full `eopkg` support in generic `pkg.Manager` interface
 - [x] Add basic ISO9660 support once again
 - [x] Add complete Legacy Boot bootloader support for `isolinux`
 - [ ] Remove repo definition from `.packages` and place in `.spin`, similar to `solbuild`.
 - [ ] Enhance bootloader support for UEFI
 - [ ] Build (successfully!) an existing Solus image specification
 - [ ] Construct specifications for our chroot builder images
 - [ ] Add support for VM/Container images

Supported Medium Types (WIP)
----------------------------

**LiveOS**

A LiveOS image is an `ISO9660` image containing a live operating system. This is the `dracut` LiveOS image type, currently used by `Solus`, `Fedora`, available in `Gentoo` and potentially others.

By default a *hybrid* ISO is created, that is an El Torito bootable image that may be booted in either an optical drive or on removal media such as a USB thumb drive. This image will use (currently) `isolinux` for the bootloader. As the project is further implemented, support will be added for `UEFI` booting too.

License
-------

Copyright © 2016 Solus Project

`USpin` is available under the terms of the Apache-2.0 license
