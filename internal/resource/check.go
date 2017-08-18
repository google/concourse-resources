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
)

type CheckRequest interface {
	Decode(source interface{}, version interface{}) error
	AddResponseVersion(version interface{})
}

type checkRequest struct {
	rawSource        json.RawMessage
	rawVersion       json.RawMessage
	responseVersions []interface{}
}

func (req checkRequest) Decode(source interface{}, version interface{}) error {
	err := json.Unmarshal(req.rawSource, source)
	if err != nil {
		return fmt.Errorf("error decoding source: %v", err)
	}

	if len(req.rawVersion) > 0 {
		err = json.Unmarshal(req.rawVersion, version)
		if err != nil {
			return fmt.Errorf("error decoding version: %v", err)
		}
	}

	return nil
}

func (req *checkRequest) AddResponseVersion(version interface{}) {
	req.responseVersions = append(req.responseVersions, version)
}

type CheckFunc func(req CheckRequest) error

func RunCheck(reqReader io.Reader, respWriter io.Writer, checkFunc CheckFunc) error {
	rawReq, err := readRawRequest(reqReader)
	if err != nil {
		return err
	}

	req := checkRequest{
		rawSource:        rawReq.Source,
		rawVersion:       rawReq.Version,
		responseVersions: []interface{}{},
	}

	err = checkFunc(&req)
	if err != nil {
		return err
	}

	return writeResponse(respWriter, req.responseVersions)
}
