package resource

import (
	"encoding/json"
	"fmt"
	"io"
)

type CheckRequest interface {
	Decode(source interface{}, version interface{}) error
	AddResponseVersion(version interface{})
}

type checkRequest struct {
	rawSource        json.RawMessage
	rawVersion       json.RawMessage
	responseVersions []interface{}
}

func (req checkRequest) Decode(source interface{}, version interface{}) error {
	err := json.Unmarshal(req.rawSource, source)
	if err != nil {
		return fmt.Errorf("error decoding source: %v", err)
	}

	err = json.Unmarshal(req.rawVersion, version)
	if err != nil {
		return fmt.Errorf("error decoding version: %v", err)
	}

	return nil
}

func (req *checkRequest) AddResponseVersion(version interface{}) {
	req.responseVersions = append(req.responseVersions, version)
}

type CheckFunc func(req CheckRequest) error

func RunCheck(reqReader io.Reader, respWriter io.Writer, checkFunc CheckFunc) error {
	rawReq, err := readRawRequest(reqReader)
	if err != nil {
		return err
	}

	req := checkRequest{
		rawSource:        rawReq.Source,
		rawVersion:       rawReq.Version,
		responseVersions: []interface{}{},
	}

	err = checkFunc(&req)
	if err != nil {
		return err
	}

	return writeResponse(respWriter, req.responseVersions)
}
