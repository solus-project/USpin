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

package spec

import (
	"reflect"
)

// OpStack contains all top level blocks
type OpStack struct {
	Blocks []*OpSet // A set of type-similar OpSet's
}

// OpSet has a given type of operations that it supports
type OpSet struct {
	Ops []Operation // Slice of operations of the same type
}

// Operation is the base Op type
type Operation interface {
	Compatible(Operation) bool
}

// An OpRepo is an operation to enable a repository on a target
type OpRepo struct {
	Operation
	RepoName string // Name to give the repository
	RepoURI  string // URI of the repository in question
}

// Compatible will always return false as OpRepo cannot be stacked
func (o *OpRepo) Compatible(o2 Operation) bool {
	return false
}

// An OpGroup is an operation to install a group or component
type OpGroup struct {
	Operation
	GroupName    string // Name of this group or component
	IgnoreSafety bool   // Whether to bypass dependency safety checks
}

// Compatible determines if two OpGroup's are compatible with one another
func (o *OpGroup) Compatible(o2 Operation) bool {
	if reflect.TypeOf(o) != reflect.TypeOf(o2) {
		return false
	}
	if o2.(*OpGroup).IgnoreSafety != o.IgnoreSafety {
		return false
	}
	return true
}

// An OpPackage is an operation to install a given package
type OpPackage struct {
	Operation
	Name         string // Name of the package to install
	IgnoreSafety bool   // Whether to bypass dependency safety checks
}

// Compatible determines if two OpPackage's are compatible with one another
func (o *OpPackage) Compatible(o2 Operation) bool {
	if reflect.TypeOf(o) != reflect.TypeOf(o2) {
		return false
	}
	if o2.(*OpPackage).IgnoreSafety != o.IgnoreSafety {
		return false
	}
	return true
}
