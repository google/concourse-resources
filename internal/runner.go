package internal

import (
	"encoding/json"
	"io"
	"reflect"
)

var resourceContextPtrType = reflect.TypeOf(&ResourceContext{})


type CheckFunc interface{}

const CheckFuncPattern = "func(Source, Version) []Version"

func validateCheckFunc(checkFunc CheckFunc) (funcValue reflect.Value, sourceType, versionType reflect.Type) {
	funcValue = reflect.ValueOf(checkFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func ||
		funcType.IsVariadic() ||
		funcType.NumIn() != 2 ||
		funcType.In(0).Kind() != reflect.Struct ||
		funcType.In(1).Kind() != reflect.Struct ||
		funcType.NumOut() != 1 ||
		funcType.Out(0).Kind() != reflect.Slice ||
		funcType.Out(0).Elem() != funcType.In(1) {
		panic("checkFunc must have signature like: " + CheckFuncPattern)
	}
	sourceType, versionType = funcType.In(0), funcType.In(1)
	return
}

func RunCheck(reqReader io.Reader, respWriter io.Writer, checkFunc CheckFunc) error {
	funcValue, sourceType, versionType := validateCheckFunc(checkFunc)

	req, err := readRequest(reqReader)
	if err != nil {
		return err
	}

	sourceValue, err := unmarshalValue(req.Source, sourceType)
	if err != nil {
		return err
	}

	versionValue, err := unmarshalValue(req.Version, versionType)
	if err != nil {
		return err
	}

	versions := call(funcValue, sourceValue, versionValue)
	return json.NewEncoder(respWriter).Encode(versions)
}

type InFunc interface{}

const InFuncSignature = "func(*ResourceContext, Source, Version, Params) Version"

func validateInFunc(inFunc InFunc) (funcValue reflect.Value, sourceType, versionType, paramsType reflect.Type) {
	funcValue = reflect.ValueOf(inFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func ||
		funcType.IsVariadic() ||
		funcType.NumIn() != 4 ||
		funcType.In(0) != resourceContextPtrType ||
		funcType.In(1).Kind() != reflect.Struct ||
		funcType.In(2).Kind() != reflect.Struct ||
		funcType.In(3).Kind() != reflect.Struct ||
		funcType.NumOut() != 1 ||
		funcType.Out(0) != funcType.In(2) {
		panic("inFunc must have signature like " + InFuncSignature)
	}
	sourceType, versionType, paramsType = funcType.In(1), funcType.In(2), funcType.In(3)
	return
}

func RunIn(targetDir string, reqReader io.Reader, respWriter io.Writer, inFunc InFunc) error {
	funcValue, sourceType, versionType, paramsType := validateInFunc(inFunc)

	req, err := readRequest(reqReader)
	if err != nil {
		return err
	}

	sourceValue, err := unmarshalValue(req.Source, sourceType)
	if err != nil {
		return err
	}

	versionValue, err := unmarshalValue(req.Version, versionType)
	if err != nil {
		return err
	}

	paramsValue, err := unmarshalValue(req.Params, paramsType)
	if err != nil {
		return err
	}

	resourceContext := ResourceContext{TargetDir: targetDir}
	version := call(funcValue, reflect.ValueOf(&resourceContext), sourceValue, versionValue, paramsValue)
	return json.NewEncoder(respWriter).Encode(ResourceResponse{
		Version: version,
		Metadata: resourceContext.Metadata,
	})
}

type OutFunc interface{}

const OutFuncSignature = "func(*ResourceContext, Source, Params) Version"

func validateOutFunc(outFunc OutFunc) (funcValue reflect.Value, sourceType, paramsType reflect.Type) {
	funcValue = reflect.ValueOf(outFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func ||
		funcType.IsVariadic() ||
		funcType.NumIn() != 3 ||
		funcType.In(0) != resourceContextPtrType ||
		funcType.In(1).Kind() != reflect.Struct ||
		funcType.In(2).Kind() != reflect.Struct ||
		funcType.NumOut() != 1 ||
		funcType.Out(0).Kind() != reflect.Struct {
		panic("outFunc must have signature like " + OutFuncSignature)
	}
	sourceType, paramsType = funcType.In(1), funcType.In(2)
	return
}	

func RunOut(targetDir string, reqReader io.Reader, respWriter io.Writer, outFunc OutFunc) error {
	funcValue, sourceType, paramsType := validateOutFunc(outFunc)

	req, err := readRequest(reqReader)
	if err != nil {
		return err
	}

	sourceValue, err := unmarshalValue(req.Source, sourceType)
	if err != nil {
		return err
	}

	paramsValue, err := unmarshalValue(req.Params, paramsType)
	if err != nil {
		return err
	}

	resourceContext := ResourceContext{TargetDir: targetDir}
	version := call(funcValue, reflect.ValueOf(&resourceContext), sourceValue, paramsValue)
	return json.NewEncoder(respWriter).Encode(ResourceResponse{
		Version: version,
		Metadata: resourceContext.Metadata,
	})
}

func readRequest(reqReader io.Reader) (req request, err error) {
	err = json.NewDecoder(reqReader).Decode(&req)
	return
}

func unmarshalValue(field json.RawMessage, valueType reflect.Type) (val reflect.Value, err error) {
	val = reflect.New(valueType)
	err = json.Unmarshal(field, val.Interface())
	if err == nil {
		val = reflect.Indirect(val)
	}
	return
}

func call(funcValue reflect.Value, in ...reflect.Value) interface{} {
	return funcValue.Call(in)[0].Interface()	
}
