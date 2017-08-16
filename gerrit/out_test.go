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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/google/concourse-resources/internal"
)

var (
	testOutVersion = Version{
		ChangeId: "outChange",
		Revision: "outRev",
	}
)

func testOut(t *testing.T, src Source, params outParams) Version {
	src.Url = testGerritUrl

	repoDir, err := ioutil.TempDir(testTempDir, "repo")
	if err != nil {
		panic(err)
	}

	err = testOutVersion.WriteToFile(filepath.Join(repoDir, gerritVersionFilename))
	if err != nil {
		panic(err)
	}
	params.Repository = filepath.Base(repoDir)


	src.Url = testGerritUrl
	req := testRequest{Source: src, Params: params}
	var resp testResourceResponse
	assert.NoError(t, internal.TestOutFunc(t, req, &resp, testTempDir, out))
	return resp.Version
}

func TestOutVersion(t *testing.T) {
	testOut(t, Source{}, outParams{})
	assert.Equal(t, "outChange", testGerritLastChangeId)
	assert.Equal(t, "outRev", testGerritLastRevision)
}

func TestOutMessage(t *testing.T) {
	testOut(t, Source{}, outParams{Message: "foo bar"})
	assert.Equal(t, "foo bar", testGerritLastReviewInput.Message)
}

func TestOutMessageFile(t *testing.T) {
	err := ioutil.WriteFile(
		filepath.Join(testTempDir, "message.txt"),
		[]byte("file msg"), 0600)
	assert.NoError(t, err)

	testOut(t, Source{}, outParams{MessageFile: "message.txt"})
	assert.Equal(t, "file msg", testGerritLastReviewInput.Message)
}

func TestOutLabels(t *testing.T) {
	testOut(t, Source{}, outParams{Labels: map[string]int{
		"Code-Review": 1,
		"Verified":    -1,
	}})
	assert.Equal(t, 1, testGerritLastReviewInput.Labels["Code-Review"])
	assert.Equal(t, -1, testGerritLastReviewInput.Labels["Verified"])
}
