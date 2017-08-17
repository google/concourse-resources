package internal

type testSource struct {
	Src string `json:"src"`
}

type testVersion struct {
	Ver int `json:"ver"`
}

type testParams struct {
	Param bool `json:"param"`
}

type testRequest struct {
	Source  testSource  `json:"source,omitempty"`
	Version testVersion `json:"version,omitempty"`
	Params  testParams  `json:"params,omitempty"`
}

var testRequestData = testRequest{
	Source:  testSource{Src: "src.go"},
	Version: testVersion{Ver: 1},
	Params:  testParams{Param: true},
}
