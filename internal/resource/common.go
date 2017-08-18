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
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ResourceRequest interface {
	TargetDir() string
	ChdirTargetDir() error

	SetResponseVersion(version interface{})
	AddResponseMetadata(key, value string)
}

type resourceResponse struct {
	Version  interface{}     `json:"version"`
	Metadata []MetadataField `json:"metadata,omitempty"`
}

type resourceRequest struct {
	targetDir string
	response  resourceResponse
}

func (req resourceRequest) TargetDir() string {
	return req.targetDir
}

func (req resourceRequest) ChdirTargetDir() error {
	return os.Chdir(req.targetDir)
}

func (req *resourceRequest) SetResponseVersion(version interface{}) {
	req.response.Version = version
}

func (req *resourceRequest) AddResponseMetadata(name string, value string) {
	req.response.Metadata = append(
		req.response.Metadata,
		MetadataField{Name: name, Value: value})
}

type MetadataField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type rawRequest struct {
	Source  json.RawMessage `json:"source"`
	Version json.RawMessage `json:"version"`
	Params  json.RawMessage `json:"params"`
}

func readRawRequest(reqReader io.Reader) (req rawRequest, err error) {
	err = json.NewDecoder(reqReader).Decode(&req)
	if err != nil {
		err = fmt.Errorf("error reading request: %v", err)
	}
	return
}

func writeResponse(respWriter io.Writer, resp interface{}) error {
	err := json.NewEncoder(respWriter).Encode(resp)
	if err != nil {
		err = fmt.Errorf("error writing response: %v", err)
	}
	return err
}
