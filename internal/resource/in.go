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

type InRequest interface {
	ResourceRequest
	Decode(source interface{}, version interface{}, params interface{}) error
}

type inRequest struct {
	resourceRequest
	rawSource  json.RawMessage
	rawVersion json.RawMessage
	rawParams  json.RawMessage
}

func (req *inRequest) Decode(source interface{}, version interface{}, params interface{}) error {
	err := json.Unmarshal(req.rawSource, source)
	if err != nil {
		return fmt.Errorf("error decoding source: %v", err)
	}

	err = json.Unmarshal(req.rawVersion, version)
	if err != nil {
		return fmt.Errorf("error decoding version: %v", err)
	}

	// Most often we just return the requested version verbatim
	if req.response.Version == nil {
		req.response.Version = version
	}

	if params != nil && len(req.rawParams) > 0 {
		err = json.Unmarshal(req.rawParams, params)
		if err != nil {
			return fmt.Errorf("error decoding params: %v", err)
		}
	}

	return nil
}

type InFunc func(req InRequest) error

func RunIn(reqReader io.Reader, respWriter io.Writer, targetDir string, inFunc InFunc) error {
	rawReq, err := readRawRequest(reqReader)
	if err != nil {
		return err
	}

	req := inRequest{
		resourceRequest: resourceRequest{targetDir: targetDir},
		rawSource:       rawReq.Source,
		rawVersion:      rawReq.Version,
		rawParams:       rawReq.Params,
	}

	err = inFunc(&req)
	if err != nil {
		return err
	}

	return writeResponse(respWriter, req.response)
}
