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

package pkg

import (
	"errors"
)

// PackageManager is the type of pkg.Manager that is to be used for packaging
// operations. This type is built-in as there is only a limited set of supported
// package managers at any time.
type PackageManager string

const (
	// PackageManagerEopkg is the package manager used within Solus
	PackageManagerEopkg PackageManager = "eopkg"
)

var (
	// ErrNotYetImplemented is a placeholder until eopkg implementation is done
	ErrNotYetImplemented = errors.New("Not yet implemented!")
)
