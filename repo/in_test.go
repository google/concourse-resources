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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/concourse-resources/internal/resource"
)

var (
	testRepoDir string
)

type testInResponse struct {
	Version `json:"version"`
	Metadata []resource.MetadataField `json:"metadata"`
}

func testIn(t *testing.T, src Source, ver Version) (Version, []resource.MetadataField) {
	var err error
	testRepoDir, err = ioutil.TempDir(testTempDir, "repo")
	assert.NoError(t, err)

	// Recreate `repo init`'s creation of .repo/manifest.xml
	assert.NoError(t, os.Mkdir(filepath.Join(testRepoDir, ".repo"), 0755))
	_, err = os.Create(filepath.Join(testRepoDir, ".repo", "manifest.xml"))
	assert.NoError(t, err)

	testLastRepoInitArgs = nil
	testLastRepoSyncArgs = nil
	src.ManifestUrl = testManifestUrl
	req := testRequest{Source: src, Version: ver}
	var resp testInResponse
	assert.NoError(t, resource.TestInFunc(t, req, &resp, testRepoDir, in))
	return resp.Version, resp.Metadata
}

func TestInVersion(t *testing.T) {
	ver := Version{Manifest: "<xml>"}
	outVer, _ := testIn(t, Source{}, ver)
	assert.Equal(t, ver, outVer)
}

func TestInRepoInit(t *testing.T) {
	testIn(t, Source{}, Version{})
	assert.Contains(t, testLastRepoInitArgs, "--manifest-url=http://fake.com/manifest")
	assert.Contains(t, testLastRepoInitArgs, "--depth=1")
}

func TestInRepoSync(t *testing.T) {
	testIn(t, Source{}, Version{})
	assert.Equal(t, testLastRepoSyncArgs[0], "sync")
}

func TestInRepoSyncManifestFile(t *testing.T) {
	ver := Version{Manifest: "<manifest-xml>"}
	testIn(t, Source{}, ver)
	manifest, err := ioutil.ReadFile(filepath.Join(testRepoDir, ".repo", "manifest.xml"))
	assert.NoError(t, err)
	assert.EqualValues(t, "<manifest-xml>", manifest)
}
