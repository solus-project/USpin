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

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// ImageSpecParser does the heavy lifting of parsing a .spin file to pull all
// relevant stack operations from it.
type ImageSpecParser struct {
	CommentCharacter   string
	RepoSplitCharacter string
	SafetyCharacter    string
	GroupCharacter     string

	Stack *OpStack

	curSet *OpSet
}

// NewParser will return a new parser for the image specification file
func NewParser() *ImageSpecParser {
	return &ImageSpecParser{
		CommentCharacter:   "#",
		RepoSplitCharacter: "=",
		SafetyCharacter:    "~",
		GroupCharacter:     "@",
		Stack:              &OpStack{},
	}
}

// pushOperation will push an operation to the last set we have
func (i *ImageSpecParser) pushOperation(op Operation) {
	if i.curSet == nil {
		i.curSet = &OpSet{}
	}
	i.curSet.Ops = append(i.curSet.Ops, op)
}

// Parse will attempt to parse the given image speicifcation file at the given
// path, and will return an error if this fails.
func (i *ImageSpecParser) Parse(path string) error {
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
				RepoName: fields[0],
				RepoURI:  value,
			}
			// TODO: Check we're adding to the right stack or need another one!
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
		fmt.Fprintf(os.Stderr, "Line: (group? %v ignoreSafety? %v) %s\n", isGroup, ignoreSafety, line)
	}

	return errors.New("Not yet implemented!")
}
