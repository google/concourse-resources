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
	assert.NoError(t, resource.TestOutFunc(t, req, &resp, testTempDir, out))
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

func TestOutMessageWithBuildId(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("BUILD_ID", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_ID}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithBuildName(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("BUILD_NAME", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_NAME}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithBuildJobName(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("BUILD_JOB_NAME", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_JOB_NAME}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithBuildPipelineName(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("BUILD_PIPELINE_NAME", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_PIPELINE_NAME}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithBuildTeamName(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("BUILD_TEAM_NAME", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_TEAM_NAME}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithATCExternalUrl(t *testing.T) {
	// Test Data
	environmentValue := "1"
	os.Setenv("ATC_EXTERNAL_URL", environmentValue)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${ATC_EXTERNAL_URL}"})

	// Verify
	assert.Equal(t, "foo bar 1", testGerritLastReviewInput.Message)
}

func TestOutMessageWithAllVariables(t *testing.T) {
	// Test Data
	buildId := "1"
	os.Setenv("BUILD_ID", buildId)
	buildName := "2"
	os.Setenv("BUILD_NAME", buildName)
	buildJobName := "3"
	os.Setenv("BUILD_JOB_NAME", buildJobName)
	buildPipelineName := "4"
	os.Setenv("BUILD_PIPELINE_NAME", buildPipelineName)
	buildTeamName := "5"
	os.Setenv("BUILD_TEAM_NAME", buildTeamName)
	atcExternalUrl := "6"
	os.Setenv("ATC_EXTERNAL_URL", atcExternalUrl)

	// Execute
	testOut(t, Source{}, outParams{Message: "foo bar ${BUILD_ID} ${BUILD_NAME} ${BUILD_JOB_NAME} ${BUILD_PIPELINE_NAME} ${BUILD_TEAM_NAME} ${ATC_EXTERNAL_URL}"})

	// Verify
	assert.Equal(t, "foo bar 1 2 3 4 5 6", testGerritLastReviewInput.Message)
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
