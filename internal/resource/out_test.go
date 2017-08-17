package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunOut(t *testing.T) {
	responseVersion := testVersion{}
	response := resourceResponse{Version: &responseVersion}
	assert.NoError(t, TestOutFunc(t, testRequestData, &response, "target", func(req OutRequest) error {
		var src testSource
		var param testParams
		assert.NoError(t, req.Decode(&src, &param))
		assert.Equal(t, "src.go", src.Src)
		assert.Equal(t, true, param.Param)
		assert.Equal(t, "target", req.TargetDir())
		req.SetResponseVersion(testVersion{Ver: 1})
		req.AddResponseMetadata("meta", "data")
		req.AddResponseMetadata("more", "meta")
		return nil
	}))
	assert.Equal(t, testVersion{Ver: 1}, responseVersion)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "meta", Value: "data"},
		MetadataField{Key: "more", Value: "meta"},
	}, response.Metadata)
}

func TestRunOutNoParams(t *testing.T) {
	assert.NoError(t, TestOutFunc(t, testRequestData, nil, "", func(req OutRequest) error {
		var src testSource
		assert.NoError(t, req.Decode(&src, nil))
		return nil
	}))
}
