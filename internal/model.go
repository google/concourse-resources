package internal

import (
	"encoding/json"
	"os"
)

type MetadataField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ResourceContext struct {
	TargetDir string
	Metadata  []MetadataField
}

func (rc ResourceContext) ChdirTarget() error {
	return os.Chdir(rc.TargetDir)
}

func (rc *ResourceContext) AddMetadata(key string, value string) {
	rc.Metadata = append(rc.Metadata, MetadataField{Key: key, Value: value})
}

type request struct {
	Source  json.RawMessage `json:"source"`
	Version json.RawMessage `json:"version"`
	Params  json.RawMessage `json:"params"`
}

type ResourceResponse struct {
	Version interface{} `json:"version"`
	Metadata []MetadataField `json:"metadata,omitempty"`
}

