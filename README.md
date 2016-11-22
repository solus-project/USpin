solspin
-----------

[![Report](https://goreportcard.com/badge/github.com/solus-project/solspin)](https://goreportcard.com/report/github.com/solus-project/solspin) [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Solus image creation utilities. Intended to succeed the existing `solus-image-creator.py` script with something a bit more robust that can construct multiple image types.

Currently the existing image creator can only construct a simple `ISO9660` image, however Solus also makes use of chrootable base images for the `overlayfs` system employed in `evobuild`.

License
-------

Copyright © 2016 Solus Project

`solspin` is available under the terms of the Apache-2.0 license
