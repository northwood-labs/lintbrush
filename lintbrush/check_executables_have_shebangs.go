// Copyright 2024, Northwood Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lintbrush

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

func CheckExecutablesHaveShebangs(run *sarif.Run, args []string) error {
	var errs *multierror.Error

	for i := range args {
		filename := args[i]

		fi, err := os.Lstat(filename)
		if err != nil {
			continue
		}

		// Make sure this file is added to the artifacts list at least once.
		CustomAddDistinctArtifact(run, fi)

		// Mask 0111 is the execute bit for ANY of owner, group, or user.
		if fi.Mode().Perm()&0111 != 0 {
			fp, err := os.Open(filename)
			if err != nil {
				errs = multierror.Append(
					errs,
					fmt.Errorf("could not open `%s`: %w", filename, err),
				)
			}

			defer fp.Close()

			reader := bufio.NewReader(fp)

			firstBytes, err := reader.Peek(2)
			if err != nil {
				errs = multierror.Append(
					errs,
					fmt.Errorf(
						"[CHECK_EXECUTABLES_HAVE_SHEBANGS] could not read the first 2 bytes of `%s`: %w",
						filename,
						err,
					),
				)
			}

			if string(firstBytes) != "#!" {
				return fmt.Errorf(
					"[CHECK_EXECUTABLES_HAVE_SHEBANGS] file %s is executable, but begins with `%s` instead of `#!`.",
					filename,
					strings.Replace(string(firstBytes), "\n", "\\n", 1),
				)
			}
		}
	}

	return errs.ErrorOrNil()
}
