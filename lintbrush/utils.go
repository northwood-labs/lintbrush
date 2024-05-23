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
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/go-enry/go-enry/v2"
	"github.com/owenrumney/go-sarif/v2/sarif"
)

func CustomAddDistinctArtifact(run *sarif.Run, fileinfo fs.FileInfo) *sarif.Artifact {
	for _, artifact := range run.Artifacts {
		if *artifact.Location.URI == fileinfo.Name() {
			return artifact
		}
	}

	hashesMap := make(map[string]string)

	// Only provide the sha256 hash if we can calculate it without errors.
	if hash, err := getSha256Hash(fileinfo); err == nil {
		hashesMap["sha-256"] = hash
	}

	a := &sarif.Artifact{
		Length: int(fileinfo.Size()),
		Hashes: hashesMap,
	}
	a.WithLocation(sarif.NewSimpleArtifactLocation(fileinfo.Name()))

	// Only provide the programming language if we can determine it without errors.
	if fd, err := os.Open(fileinfo.Name()); err == nil {
		reader := bufio.NewReader(fd)
		if byteSet, err := reader.Peek(2048); err == nil {
			lang, safe := enry.GetLanguageByContent(
				filepath.Base(fileinfo.Name()),
				byteSet,
			)

			if safe {
				a.WithSourceLanguage(lang)
			}
		}
	}

	run.Artifacts = append(run.Artifacts, a)

	return a
}

func getSha256Hash(fileinfo fs.FileInfo) (string, error) {
	h := sha256.New()

	f, err := os.Open(fileinfo.Name())
	if err != nil {
		return "", fmt.Errorf("could not open `%s`: %w", fileinfo.Name(), err)
	}

	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("could calculate hash for `%s`: %w", fileinfo.Name(), err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
