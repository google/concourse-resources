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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptionsUnmarshal(t *testing.T) {
	var opts options
	assert.NoError(t, json.Unmarshal([]byte(`{"a": "val", "b": true}`), &opts))
	assert.Equal(t, "val", opts["a"])
	assert.Equal(t, true, opts["b"])
}

func TestOptionsUnmarshalBadType(t *testing.T) {
	var opts options
	err := json.Unmarshal([]byte(`{"bad": null}`), &opts)
	assert.EqualError(t, err, `invalid option type <nil> for option "bad"`)
}

func TestOptionsMerge(t *testing.T) {
	opts := options{"a": "A", "b": "B"}
	opts.merge(options{"b": "X", "c": "C"})
	assert.Equal(t, options{"a": "A", "b": "X", "c": "C"}, opts)
}

func TestOptionsArgs(t *testing.T) {
	args := options{"flag": true, "noflag": false, "num": float64(4), "opt": "val"}.args()
	assert.Equal(t, []string{"--flag", "--num=4", "--opt=val"}, args)
}
