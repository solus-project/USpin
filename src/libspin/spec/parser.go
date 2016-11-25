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

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Parser does the heavy lifting of parsing a .spin file to pull all
// relevant stack operations from it.
type Parser struct {
	CommentCharacter   string // If a line starts with this character, it is ignored. Defaults to '#'
	RepoSplitCharacter string // Character to denote a repo definition. Defaults to '='
	SafetyCharacter    string // Character to indicate ignoreSafety. Defaults to '~'
	GroupCharacter     string // Character to indicate a group or component. Defaults to '@'

	Stack *OpStack // The parsed stack so far

	curSet *OpSet
}

// NewParser will return a new parser for the image specification file
func NewParser() *Parser {
	return &Parser{
		CommentCharacter:   "#",
		RepoSplitCharacter: "=",
		SafetyCharacter:    "~",
		GroupCharacter:     "@",
		Stack:              &OpStack{},
	}
}

// pushOperation will push an operation to the last set we have
func (i *Parser) pushOperation(op Operation) {
	if i.curSet == nil {
		i.curSet = &OpSet{}
	}
	// Just insert the operation, it's the first guy
	if len(i.curSet.Ops) == 0 {
		goto insertOp
	}
	// Type mismatch, begin a new stack
	if !op.Compatible(i.curSet.Ops[0]) {
		i.Stack.Blocks = append(i.Stack.Blocks, i.curSet)
		i.curSet = &OpSet{}
	}

insertOp:
	i.curSet.Ops = append(i.curSet.Ops, op)
}

// Parse will attempt to parse the given image speicifcation file at the given
// path, and will return an error if this fails.
func (i *Parser) Parse(path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	sc := bufio.NewScanner(fi)

	lineno := 0

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		lineno++

		// ~ character ignores safety.
		ignoreSafety := false
		isGroup := false

		if line == "" {
			continue
		}

		// Check for single line comments
		if strings.HasPrefix(line, i.CommentCharacter) {
			continue
		}

		// Check if this is a repo
		if strings.Contains(line, i.RepoSplitCharacter) {
			fields := strings.Split(line, "=")
			value := strings.TrimSpace(strings.Join(fields[1:], "="))
			if value == "" {
				return fmt.Errorf("Missing value for repo declaration '%v' on line '%v'\n", fields[0], lineno)
			}
			op := &OpRepo{
				RepoName: strings.TrimSpace(fields[0]),
				RepoURI:  value,
			}
			i.pushOperation(op)
			continue
		}

		// Check if safety is disabled
		if strings.HasPrefix(line, i.SafetyCharacter) {
			ignoreSafety = true
			line = line[len(i.SafetyCharacter):]
		}

		// Check if its a group or not
		if strings.HasPrefix(line, i.GroupCharacter) {
			isGroup = true
			line = line[len(i.GroupCharacter):]
		}

		var op Operation

		// Add the operation to the stack
		if isGroup {
			op = &OpGroup{
				GroupName:    line,
				IgnoreSafety: ignoreSafety,
			}
		} else {
			op = &OpPackage{
				Name:         line,
				IgnoreSafety: ignoreSafety,
			}
		}
		i.pushOperation(op)
	}

	i.Stack.Blocks = append(i.Stack.Blocks, i.curSet)
	i.curSet = nil

	return nil
}
