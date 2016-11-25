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

// Package spec contains the specification format and parser for .spin package files
//
// Format
//
// A spin package file is a very simple text file, which supports comments and
// minimal command style.
// Any line that begins with a '#', or is blank after trimming, is ignored.
//
// Repo Lines
//
// A repository is defined as having a key and value. If a line is encountered
// using the split delimiter '=', the left side is assumed to be the name of the
// repository, and the right side is assumed to be the URI of this repository.
//      RepoName = http://example.com/eopkg-index.xml.xz
//
// Group lines
//
// A line beginning with the group character '@' is interpreted as a request
// to install the named group, whose name shall be the line minus the '@'.
// Given the input:
//      @system.base
// The component named "system.base" would be installed.
//
// Package lines
//
// Any non blank line neither qualifying as a repo or group line is interpreted
// as a package installation.
//
// Control Characters
//
// An additional character, '~', may be used by implementations to control the
// 'IgnoreSafety' parameter of Package & Group install lines. Depending on the
// implementation, this will bypass dependency safety checks in order to break
// a cyclical dependency to inject a group or package before other dependencies
// are met, such as for baselayout style packages.
//
// This control character must be the first character in the sequence.
package spec
