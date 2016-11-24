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

package main

import (
	"fmt"
	_ "libimage"
	"libspin"
	"os"
)

func printUsage(exitCode int) {
	var fd *os.File
	if exitCode == 0 {
		fd = os.Stdout
	} else {
		fd = os.Stderr
	}

	fmt.Fprintf(fd, "%s [image.spin]\n", os.Args[0])
	os.Exit(exitCode)
}

func main() {
	if len(os.Args) < 2 {
		printUsage(1)
	}

	if _, err := libspin.NewImageSpec(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Not yet implemented\n")
}
