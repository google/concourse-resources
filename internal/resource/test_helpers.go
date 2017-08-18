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
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFunc(t *testing.T, request interface{}, response interface{}, checkFunc CheckFunc) error {
	return testRunner(t, RunCheck, request, response, checkFunc)
}

func TestInFunc(t *testing.T, request interface{}, response interface{}, targetDir string, inFunc InFunc) error {
	return testRunner(t, RunIn, request, response, targetDir, inFunc)
}

func TestOutFunc(t *testing.T, request interface{}, response interface{}, targetDir string, outFunc OutFunc) error {
	return testRunner(t, RunOut, request, response, targetDir, outFunc)
}

func testRunner(t *testing.T, runner interface{}, req interface{}, resp interface{}, args ...interface{}) error {
	requestBuf := new(bytes.Buffer)
	assert.NoError(t, json.NewEncoder(requestBuf).Encode(req))

	if t.Failed() {
		t.FailNow()
	}

	responseBuf := new(bytes.Buffer)

	argVals := []reflect.Value{
		reflect.ValueOf(requestBuf),
		reflect.ValueOf(responseBuf),
	}
	for _, arg := range args {
		argVals = append(argVals, reflect.ValueOf(arg))
	}

	results := reflect.ValueOf(runner).Call(argVals)

	if t.Failed() {
		t.FailNow()
	}

	if resp != nil {
		assert.NoError(t, json.NewDecoder(responseBuf).Decode(resp))
	}

	assert.Len(t, results, 1)

	if results[0].IsNil() {
		return nil
	} else {
		return results[0].Interface().(error)
	}
}
