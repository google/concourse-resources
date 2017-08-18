// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	testManifestUrl = "http://fake.com/manifest"
)

var (
	testTempDir string
)

type testRequest struct {
	Source  `json:"source"`
	Version `json:"version"`
}

func TestMain(m *testing.M) {
	// Run a separate func so defers run before Exit.
	os.Exit(func() int {
		testingRepo = true

		var err error
		testTempDir, err = ioutil.TempDir("", "concourse-repo-test")
		if err != nil {
			panic(err)
		}
		defer os.RemoveAll(testTempDir)

		return m.Run()
	}())
}
