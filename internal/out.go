package internal

import (
	"encoding/json"
	"fmt"
	"io"
)

type OutRequest interface {
	ResourceRequest
	Decode(source interface{}, params interface{}) error
}

type outRequest struct {
	resourceRequest
	rawSource json.RawMessage
	rawParams json.RawMessage
}

func (req outRequest) Decode(source interface{}, params interface{}) error {
	err := json.Unmarshal(req.rawSource, source)
	if err != nil {
		return fmt.Errorf("error decoding source: %v", err)
	}

	if params != nil {
		err = json.Unmarshal(req.rawParams, params)
		if err != nil {
			return fmt.Errorf("error decoding params: %v", err)
		}
	}

	return nil
}

type OutFunc func(req OutRequest) error

func RunOut(reqReader io.Reader, respWriter io.Writer, targetDir string, outFunc OutFunc) error {
	rawReq, err := readRawRequest(reqReader)
	if err != nil {
		return err
	}

	req := outRequest{
		resourceRequest: resourceRequest{targetDir: targetDir},
		rawSource:   rawReq.Source,
		rawParams:   rawReq.Params,
	}

	err = outFunc(&req)
	if err != nil {
		return err
	}

	return writeResponse(respWriter, req.response)
}
