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
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/concourse-resources/internal/resource"
)

var (
	testCheckRunCount = 0
)

func testCheck(t *testing.T, src Source, ver Version) []Version {
	// Run each test in a separate subdir.
	testCheckRunCount++
	checkRepoDir = filepath.Join(testTempDir, fmt.Sprintf("repo%d", testCheckRunCount))

	testLastRepoInitArgs = nil
	testLastRepoSyncArgs = nil
	src.ManifestUrl = testManifestUrl
	req := testRequest{Source: src, Version: ver}
	var versions []Version
	assert.NoError(t, resource.TestCheckFunc(t, req, &versions, check))
	if t.Failed() {
		t.FailNow()
	}
	return versions
}

func TestCheckVersions(t *testing.T) {
	testCurrentVersion = Version{Manifest: "<xml>"}
	otherVersion := Version{Manifest: "<otherxml>"}

	// Request empty version, get current version
	versions := testCheck(t, Source{}, Version{})
	assert.Equal(t, []Version{testCurrentVersion}, versions)

	// Request current version, get current version
	versions = testCheck(t, Source{}, testCurrentVersion)
	assert.Equal(t, []Version{testCurrentVersion}, versions)

	// Request other version, get other and current versions
	versions = testCheck(t, Source{}, otherVersion)
	assert.Equal(t, []Version{otherVersion, testCurrentVersion}, versions)
}

func TestCheckRepoInit(t *testing.T) {
	testCheck(t, Source{}, Version{})
	assert.Contains(t, testLastRepoInitArgs, "--manifest-url=http://fake.com/manifest")
	assert.Contains(t, testLastRepoInitArgs, "--depth=1")

	// Shouldn't run repoInit on second run.
	testCheckRunCount--
	testCheck(t, Source{}, Version{})
	assert.Nil(t, testLastRepoInitArgs)
}

func TestCheckRepoSync(t *testing.T) {
	testCheck(t, Source{}, Version{})
	assert.Equal(t, testLastRepoSyncArgs[0], "sync")

	// Should run repoSync on second run.
	testCheckRunCount--
	testCheck(t, Source{}, Version{})
	assert.NotNil(t, testLastRepoSyncArgs)
}
