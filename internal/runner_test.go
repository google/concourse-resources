package internal

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Version  testVersion     `json:"version"`
	Metadata []MetadataField `json:"metadata"`
}

func TestRunCheck(t *testing.T) {
	req := testRequest{
		Source:  testSource{S: "s"},
		Version: testVersion{V: 1},
	}
	var resp []testVersion
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		assert.NoError(t, RunCheck(in, out,
			func(s testSource, v testVersion) ([]testVersion, error) {
				assert.Equal(t, req.Source, s)
				assert.Equal(t, req.Version, v)
				return []testVersion{testVersion{V: 2}}, nil
			}))
	})
	assert.Equal(t, []testVersion{testVersion{V: 2}}, resp)
}

func TestRunCheckError(t *testing.T) {
	req := testRequest{
		Source:  testSource{S: "s"},
		Version: testVersion{V: 1},
	}
	var resp []testVersion
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		assert.EqualError(t, RunCheck(in, out,
			func(s testSource, v testVersion) ([]testVersion, error) {
				return nil, errors.New("my error")
			}), "my error")
	})
}

func TestValidateCheckFunc(t *testing.T) {
	testFuncs := []CheckFunc{
		testSource{},
		func(testSource, testVersion) []testVersion { return nil },
		func(testSource, testVersion) testVersion { return testVersion{} },
		func(testSource, ...testVersion) ([]testVersion, error) { return nil, nil },
		func(testSource, testVersion, int) ([]testVersion, error) { return nil, nil },
		func(int, testVersion) ([]testVersion, error) { return nil, nil },
		func(testSource, testVersion) ([]int, error) { return nil, nil },
		func(testSource, testVersion) ([]testVersion, int) { return nil, 0 },
		func(testSource, testVersion) ([]testSource, error) { return nil, nil },
	}
	for _, testFunc := range testFuncs {
		assert.Panics(t, func() { validateCheckFunc(testFunc) })
	}
}

func TestRunIn(t *testing.T) {
	req := testRequest{
		Source:  testSource{S: "s"},
		Version: testVersion{V: 1},
		Params:  testParams{P: true},
	}
	var resp testResponse
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		assert.NoError(t, RunIn("foo", in, out,
			func(rc *ResourceContext, s testSource, v testVersion, p testParams) (testVersion, error) {
				assert.Equal(t, "foo", rc.TargetDir)
				assert.Equal(t, req.Source, s)
				assert.Equal(t, req.Version, v)
				assert.Equal(t, req.Params, p)
				rc.AddMetadata("x", "1")
				rc.AddMetadata("y", "2")
				return testVersion{V: 2}, nil
			}))
	})
	assert.Equal(t, testVersion{V: 2}, resp.Version)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "x", Value: "1"},
		MetadataField{Key: "y", Value: "2"},
	}, resp.Metadata)
}

func TestRunInError(t *testing.T) {
	req := testRequest{
		Source:  testSource{S: "s"},
		Version: testVersion{V: 1},
		Params:  testParams{P: true},
	}
	var resp testResponse
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		assert.EqualError(t, RunIn("foo", in, out,
			func(rc *ResourceContext, s testSource, v testVersion, p testParams) (testVersion, error) {
				return testVersion{}, errors.New("my error")
			}), "my error")
	})
}

func TestValidateInFunc(t *testing.T) {
	testFuncs := []InFunc{
		testSource{},
		func(*ResourceContext, testSource, testVersion, ...testParams) (testVersion, error) {
			return testVersion{}, nil
		},
		func(*ResourceContext, testSource, testVersion, int) (testVersion, error) { return testVersion{}, nil },
		func(*ResourceContext, testSource, testVersion, testParams, int) (testVersion, error) {
			return testVersion{}, nil
		},
		func(*ResourceContext, testSource, testVersion, testParams) (int, error) { return 0, nil },
		func(*ResourceContext, testSource, testVersion, testParams) (testVersion, int, error) {
			return testVersion{}, 0, nil
		},
		func(*ResourceContext, testSource, testVersion, testParams) (testParams, error) {
			return testParams{}, nil
		},
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
		assert.NoError(t, RunOut("foo", in, out,
			func(rc *ResourceContext, s testSource, p testParams) (testVersion, error) {
				assert.Equal(t, "foo", rc.TargetDir)
				assert.Equal(t, req.Source, s)
				assert.Equal(t, req.Params, p)
				rc.AddMetadata("x", "1")
				rc.AddMetadata("y", "2")
				return testVersion{V: 2}, nil
			}))
	})
	assert.Equal(t, testVersion{V: 2}, resp.Version)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "x", Value: "1"},
		MetadataField{Key: "y", Value: "2"},
	}, resp.Metadata)
}

func TestRunOutError(t *testing.T) {
	req := testRequest{
		Source: testSource{S: "s"},
		Params: testParams{P: true},
	}
	var resp testResponse
	jsonInOut(req, &resp, func(in io.Reader, out io.Writer) {
		assert.EqualError(t, RunOut("foo", in, out,
			func(rc *ResourceContext, s testSource, p testParams) (testVersion, error) {
				return testVersion{}, errors.New("my error")
			}), "my error")
	})
}

func TestValidateOutFunc(t *testing.T) {
	testFuncs := []interface{}{
		testSource{},
		func(*ResourceContext, testSource, ...testParams) (testVersion, error) { return testVersion{}, nil },
		func(*ResourceContext, testSource, int) (testVersion, error) { return testVersion{}, nil },
		func(*ResourceContext, testSource, testParams, int) (testVersion, error) { return testVersion{}, nil },
		func(*ResourceContext, testSource, testParams) (int, error) { return 0, nil },
		func(*ResourceContext, testSource, testParams) (testVersion, int, error) { return testVersion{}, 0, nil },
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
	if out.Len() > 0 {
		panicOnError(json.NewDecoder(&out).Decode(outVal))
	}
}
