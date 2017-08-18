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
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRepoInit(t *testing.T, src Source) {
	testLastRepoInitArgs = nil
	src.ManifestUrl = testManifestUrl
	assert.NoError(t, repoInit("/tmp/fake", src))
}

func testRepoSync(t *testing.T, src Source) {
	testLastRepoSyncArgs = nil
	src.ManifestUrl = testManifestUrl
	assert.NoError(t, repoSync("/tmp/fake", src))
}

func TestRepoInitManifestUrlRequired(t *testing.T) {
	assert.EqualError(t, repoInit("", Source{}), "manifest_url is required")
}

func TestRepoInitOptions(t *testing.T) {
	testRepoInit(t, Source{
		ManifestName:   "other.xml",
		ManifestBranch: "mybranch",
		Groups:         []string{"group1", "group2"},
		InitOptions: options{
			"reference": "/mirror",
			"quiet":     true,
			"depth":     false,
		},
	})

	assert.Contains(t, testLastRepoInitArgs, "--manifest-name=other.xml")
	assert.Contains(t, testLastRepoInitArgs, "--manifest-branch=mybranch")
	assert.Contains(t, testLastRepoInitArgs, "--groups=group1,group2")
	assert.Contains(t, testLastRepoInitArgs, "--quiet")
	assert.Contains(t, testLastRepoInitArgs, "--reference=/mirror")
	assert.NotContains(t, testLastRepoInitArgs, "--depth=1")
}

func TestRepoSyncOptions(t *testing.T) {
	testRepoSync(t, Source{
		SyncOptions: options{
			"jobs":            float64(99),
			"manifest-name":   "source.xml",
			"optimized-fetch": false,
		},
	})
	assert.Contains(t, testLastRepoSyncArgs, "--jobs=99")
	assert.NotContains(t, testLastRepoSyncArgs, "--optimized-fetch")
}
