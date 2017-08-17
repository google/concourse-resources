package resource

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ResourceRequest interface {
	TargetDir() string
	ChdirTargetDir() error

	SetResponseVersion(version interface{})
	AddResponseMetadata(key, value string)
}

type resourceResponse struct {
	Version  interface{}     `json:"version"`
	Metadata []MetadataField `json:"metadata,omitempty"`
}

type resourceRequest struct {
	targetDir string
	response  resourceResponse
}

func (req resourceRequest) TargetDir() string {
	return req.targetDir
}

func (req resourceRequest) ChdirTargetDir() error {
	return os.Chdir(req.targetDir)
}

func (req *resourceRequest) SetResponseVersion(version interface{}) {
	req.response.Version = version
}

func (req *resourceRequest) AddResponseMetadata(key string, value string) {
	req.response.Metadata = append(
		req.response.Metadata,
		MetadataField{Key: key, Value: value})
}

type MetadataField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type rawRequest struct {
	Source  json.RawMessage `json:"source"`
	Version json.RawMessage `json:"version"`
	Params  json.RawMessage `json:"params"`
}

func readRawRequest(reqReader io.Reader) (req rawRequest, err error) {
	err = json.NewDecoder(reqReader).Decode(&req)
	if err != nil {
		err = fmt.Errorf("error reading request: %v", err)
	}
	return
}

func writeResponse(respWriter io.Writer, resp interface{}) error {
	err := json.NewEncoder(respWriter).Encode(resp)
	if err != nil {
		err = fmt.Errorf("error writing response: %v", err)
	}
	return err
}
