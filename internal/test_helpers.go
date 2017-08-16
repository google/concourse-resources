package internal

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

	responseBuf := new(bytes.Buffer)
	argVals := []reflect.Value{
		reflect.ValueOf(requestBuf),
		reflect.ValueOf(responseBuf),
	}
	for _, arg := range args {
		argVals = append(argVals, reflect.ValueOf(arg))
	}
	results := reflect.ValueOf(runner).Call(argVals)
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
