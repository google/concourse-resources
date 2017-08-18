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

package resource

type testSource struct {
	Src string `json:"src"`
}

type testVersion struct {
	Ver int `json:"ver"`
}

type testParams struct {
	Param bool `json:"param"`
}

type testRequest struct {
	Source  testSource  `json:"source,omitempty"`
	Version testVersion `json:"version,omitempty"`
	Params  testParams  `json:"params,omitempty"`
}

var testRequestData = testRequest{
	Source:  testSource{Src: "src.go"},
	Version: testVersion{Ver: 1},
	Params:  testParams{Param: true},
}
