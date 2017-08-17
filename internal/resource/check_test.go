package resource

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCheck(t *testing.T) {
	var versions []testVersion
	assert.NoError(t, TestCheckFunc(t, testRequestData, &versions, func(req CheckRequest) error {
		var src testSource
		var ver testVersion
		assert.NoError(t, req.Decode(&src, &ver))
		assert.Equal(t, "src.go", src.Src)
		assert.Equal(t, 1, ver.Ver)
		req.AddResponseVersion(testVersion{Ver: 1})
		req.AddResponseVersion(testVersion{Ver: 2})
		return nil
	}))
	assert.Equal(t, []testVersion{
		testVersion{Ver: 1},
		testVersion{Ver: 2},
	}, versions)
}

func TestRunCheckNoResults(t *testing.T) {
	var versions json.RawMessage
	assert.NoError(t, TestCheckFunc(t, testRequestData, &versions, func(req CheckRequest) error {
		return nil
	}))
	assert.Equal(t, "[]", string(versions))
}

func TestRunCheckError(t *testing.T) {
	err := errors.New("my error")
	assert.Equal(t, err, TestCheckFunc(t, testRequestData, nil, func(req CheckRequest) error {
		return err
	}))
}
