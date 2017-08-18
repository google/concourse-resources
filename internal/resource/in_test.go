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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunIn(t *testing.T) {
	responseVersion := testVersion{}
	response := resourceResponse{Version: &responseVersion}
	assert.NoError(t, TestInFunc(t, testRequestData, &response, "target", func(req InRequest) error {
		var src testSource
		var ver testVersion
		var param testParams
		assert.NoError(t, req.Decode(&src, &ver, &param))
		assert.Equal(t, "src.go", src.Src)
		assert.Equal(t, 1, ver.Ver)
		assert.Equal(t, true, param.Param)
		assert.Equal(t, "target", req.TargetDir())
		req.SetResponseVersion(testVersion{Ver: 2})
		req.AddResponseMetadata("meta", "data")
		req.AddResponseMetadata("more", "meta")
		return nil
	}))
	assert.Equal(t, testVersion{Ver: 2}, responseVersion)
	assert.Equal(t, []MetadataField{
		MetadataField{Name: "meta", Value: "data"},
		MetadataField{Name: "more", Value: "meta"},
	}, response.Metadata)
}

func TestRunInNoSetResponseVersion(t *testing.T) {
	responseVersion := testVersion{}
	response := resourceResponse{Version: &responseVersion}
	assert.NoError(t, TestInFunc(t, testRequestData, &response, "", func(req InRequest) error {
		var src testSource
		var ver testVersion
		assert.NoError(t, req.Decode(&src, &ver, nil))
		return nil
	}))
	assert.Equal(t, testVersion{Ver: 1}, responseVersion)
}
