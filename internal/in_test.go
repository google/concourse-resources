package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunIn(t *testing.T) {
	responseVersion := testVersion{}
	response := resourceResponse{Version: &responseVersion}
	assert.NoError(t, TestInFunc(t, testRequestData, &response, "target", func(req InRequest) error {
		var src testSource
		var ver testVersion
		var param testParams
		assert.NoError(t, req.Decode(&src, &ver, &param))
		assert.Equal(t, "src.go", src.Src)
		assert.Equal(t, 1, ver.Ver)
		assert.Equal(t, true, param.Param)
		assert.Equal(t, "target", req.TargetDir())
		req.SetResponseVersion(testVersion{Ver: 2})
		req.AddResponseMetadata("meta", "data")
		req.AddResponseMetadata("more", "meta")
		return nil
	}))
	assert.Equal(t, testVersion{Ver: 2}, responseVersion)
	assert.Equal(t, []MetadataField{
		MetadataField{Key: "meta", Value: "data"},
		MetadataField{Key: "more", Value: "meta"},
	}, response.Metadata)
}

func TestRunInNoSetResponseVersion(t *testing.T) {
	responseVersion := testVersion{}
	response := resourceResponse{Version: &responseVersion}
	assert.NoError(t, TestInFunc(t, testRequestData, &response, "", func(req InRequest) error {
		var src testSource
		var ver testVersion
		assert.NoError(t, req.Decode(&src, &ver, nil))
		return nil
	}))
	assert.Equal(t, testVersion{Ver: 1}, responseVersion)
}
