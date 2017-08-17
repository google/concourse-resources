package internal

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

	if params != nil {
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
