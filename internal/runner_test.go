package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSource struct {
	S string
}

type testParams struct {
	P bool
}

type testVersion struct {
	V int
}

type testRequest struct {
	Source  testSource  `json:"source"`
	Version testVersion `json:"version"`
	Params  testParams  `json:"params"`
}

type testResponse struct {
	Version testVersion `json:"version"`
	Metadata []MetadataField `json:"metadata"`
}

func TestRunCheck(t *testing.T) {
	req := testRequest{
		Source: testSource{S: "s"},
		Version: testVersion{V: 1},
	}
	var resp []testVersion
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		RunCheck(in, out, func(s testSource, v testVersion) []testVersion {
			assert.Equal(t, req.Source, s)
			assert.Equal(t, req.Version, v)
			return []testVersion{testVersion{V: 2}}
		})
	})
	assert.Equal(t, []testVersion{testVersion{V: 2}}, resp)
}

func TestValidateCheckFunc(t *testing.T) {
	testFuncs := []CheckFunc{
		testSource{},
		func(testSource, ...testVersion) []testVersion { return nil },
		func(testSource, testVersion, int) []testVersion { return nil },
		func(int, testVersion) []testVersion { return nil },
		func(testSource, testVersion) []int { return nil },
		func(testSource, testVersion) ([]testVersion, int) { return nil, 0 },
		func(testSource, testVersion) []testSource { return nil },
	}
	for _, testFunc := range testFuncs {
		assert.Panics(t, func() { validateCheckFunc(testFunc) })
	}
}

func TestRunIn(t *testing.T) {
	req := testRequest{
		Source: testSource{S: "s"},
		Version: testVersion{V: 1},
		Params: testParams{P: true},
	}
	var resp testResponse
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		RunIn("foo", in, out, func(rc *ResourceContext, s testSource, v testVersion, p testParams) testVersion {
			assert.Equal(t, "foo", rc.TargetDir)
			assert.Equal(t, req.Source, s)
			assert.Equal(t, req.Version, v)
			assert.Equal(t, req.Params, p)
			rc.AddMetadata("x", "1")
			rc.AddMetadata("y", "2")
			return testVersion{V: 2}
		})
	})
	assert.Equal(t, testVersion{V: 2}, resp.Version)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "x", Value: "1"},
		MetadataField{Key: "y", Value: "2"},
	}, resp.Metadata)
}

func TestValidateInFunc(t *testing.T) {
	testFuncs := []InFunc{
		testSource{},
		func(*ResourceContext, testSource, testVersion, ...testParams) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, testVersion, int) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, testVersion, testParams, int) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, testVersion, testParams) int { return 0 },
		func(*ResourceContext, testSource, testVersion, testParams) (testVersion, int) { return testVersion{}, 0 },
		func(*ResourceContext, testSource, testVersion, testParams) testParams { return testParams{} },
	}
	for _, testFunc := range testFuncs {
		assert.Panics(t, func() { validateInFunc(testFunc) })
	}
}

func TestRunOut(t *testing.T) {
	req := testRequest{
		Source: testSource{S: "s"},
		Params: testParams{P: true},
	}
	var resp testResponse
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		RunOut("foo", in, out, func(rc *ResourceContext, s testSource, p testParams) testVersion {
			assert.Equal(t, "foo", rc.TargetDir)
			assert.Equal(t, req.Source, s)
			assert.Equal(t, req.Params, p)
			rc.AddMetadata("x", "1")
			rc.AddMetadata("y", "2")
			return testVersion{V: 2}
		})
	})
	assert.Equal(t, testVersion{V: 2}, resp.Version)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "x", Value: "1"},
		MetadataField{Key: "y", Value: "2"},
	}, resp.Metadata)
}

func TestValidateOutFunc(t *testing.T) {
	testFuncs := []interface{}{
		testSource{},
		func(*ResourceContext, testSource, ...testParams) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, int) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, testParams, int) testVersion { return testVersion{} },
		func(*ResourceContext, testSource, testParams) int { return 0 },
		func(*ResourceContext, testSource, testParams) (testVersion, int) { return testVersion{}, 0 },
	}
	for _, testFunc := range testFuncs {
		assert.Panics(t, func() { validateOutFunc(testFunc) })
	}
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func jsonInOut(inVal interface{}, outVal interface{}, f func(io.Reader, io.Writer)) {
	var in, out bytes.Buffer
	panicOnError(json.NewEncoder(&in).Encode(inVal))
	f(&in, &out)
	panicOnError(json.NewDecoder(&out).Decode(outVal))
}
