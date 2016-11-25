USpin
------

[![Report](https://goreportcard.com/badge/github.com/solus-project/USpin)](https://goreportcard.com/report/github.com/solus-project/USpin) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Universal image creation utilities. Used to produce bootable Linux medium in an agnostic fashion.

[![You Spin Me Right Round](https://img.youtube.com/vi/PGNiXGX2nLU/0.jpg)](https://www.youtube.com/watch?v=PGNiXGX2nLU)

Rationale
-------
Intended to succeed the existing `solus-image-creator.py` script with something a bit more robust that can construct multiple image types.

Currently the existing image creator can only construct a simple `ISO9660` image, however Solus also makes use of chrootable base images for the `overlayfs` system employed in `evobuild`.

TODO
----

 - [x] Add parser for the Solus image specification format
 - [x] Port the `Stack` implementation from old image creator
 - [x] Add config format for the main image configuration
 - [x] Add utilities for image format & creation (`dd`/`fallocate`, etc)
 - [x] Implement full `eopkg` support in generic `pkg.Manager` interface
 - [ ] Add basic ISO9660 support once again
 - [ ] Build (successfully!) an existing Solus image specification
 - [ ] Construct specifications for our chroot builder images
 - [ ] Add support for VM/Container images


License
-------

Copyright © 2016 Solus Project

`USpin` is available under the terms of the Apache-2.0 license
