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
	"fmt"
	"sort"
)

type Source struct {
	// Init options
	ManifestUrl    string   `json:"manifest_url"`
	ManifestName   string   `json:"manifest_name"`
	ManifestBranch string   `json:"manifest_branch"`
	Groups         []string `json:"groups"`
	InitOptions    options  `json:"init_options"`

	// Sync options
	SyncOptions options `json:"sync_options"`
}

type Version struct {
	Manifest string `json:"manifest"`
}

type options map[string]interface{}

func (opts *options) UnmarshalJSON(data []byte) error {
	var optsMap map[string]interface{}
	err := json.Unmarshal(data, &optsMap)
	if err != nil {
		return err
	}
	for key, val := range optsMap {
		switch val.(type) {
		case bool, string, float64:
		default:
			return fmt.Errorf("invalid option type %T for option %q", val, key)
		}
	}
	*opts = options(optsMap)
	return nil
}

func (opts options) merge(other options) {
	for key, val := range other {
		opts[key] = val
	}
}

func (opts options) args() []string {
	args := []string{}
	for key, val := range opts {
		var arg string

		switch v := val.(type) {
		case bool:
			if v {
				arg = fmt.Sprintf("--%s", key)
			} else {
				continue // Ignore args with value `false`
			}
		case string, float64:
			arg = fmt.Sprintf("--%s=%v", key, v)
		default:
			panic("options values must be bool, string, or number")
		}
		args = append(args, arg)
	}
	sort.Strings(args)
	return args
}
