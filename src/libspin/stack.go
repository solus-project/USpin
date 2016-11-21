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

package libspin

// OpStack contains all top level blocks
type OpStack struct {
	Blocks []*OpSet
}

// OpSet has a given type of operations that it supports
type OpSet struct {
	Ops []*Operation
}

// Operation is the base Op type
type Operation interface{}

// An OpRepo is an operation to enable a repository on a target
type OpRepo struct {
	Operation
	RepoName string
	RepoURI  string
}

// An OpGroup is an operation to install a group/component
type OpGroup struct {
	Operation
	GroupName string
}

// An OpPackage is an operation to install a given package
type OpPackage struct {
	Operation
	Name         string
	IgnoreSafety bool
}
