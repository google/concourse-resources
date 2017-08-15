package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

var (
	resourceContextPtrType = reflect.TypeOf(&ResourceContext{})
	errorPtr               *error
	errorType              = reflect.TypeOf(errorPtr).Elem()
)

type CheckFunc interface{}

const CheckFuncPattern = "func(Source, Version) ([]Version, error)"

func validateCheckFunc(checkFunc CheckFunc) (funcValue reflect.Value, sourceType, versionType reflect.Type) {
	funcValue = reflect.ValueOf(checkFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func ||
		funcType.IsVariadic() ||
		funcType.NumIn() != 2 ||
		funcType.In(0).Kind() != reflect.Struct ||
		funcType.In(1).Kind() != reflect.Struct ||
		funcType.NumOut() != 2 ||
		funcType.Out(0).Kind() != reflect.Slice ||
		funcType.Out(0).Elem() != funcType.In(1) ||
		funcType.Out(1) != errorType {
		panic(fmt.Sprintf("checkFunc must have signature like '%s', got '%T'", CheckFuncPattern, checkFunc))
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

	versions, err := call(funcValue, sourceValue, versionValue)
	if err != nil {
		return err
	}

	return json.NewEncoder(respWriter).Encode(versions)
}

type InFunc interface{}

const InFuncPattern = "func(*ResourceContext, Source, Version, Params) (Version, error)"

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
		funcType.NumOut() != 2 ||
		funcType.Out(0) != funcType.In(2) ||
		funcType.Out(1) != errorType {
		panic(fmt.Sprintf("inFunc must have signature like '%s', got '%T'", InFuncPattern, inFunc))
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
	version, err := call(funcValue, reflect.ValueOf(&resourceContext), sourceValue, versionValue, paramsValue)
	if err != nil {
		return err
	}

	return json.NewEncoder(respWriter).Encode(ResourceResponse{
		Version:  version,
		Metadata: resourceContext.Metadata,
	})
}

type OutFunc interface{}

const OutFuncPattern = "func(*ResourceContext, Source, Params) (Version, error)"

func validateOutFunc(outFunc OutFunc) (funcValue reflect.Value, sourceType, paramsType reflect.Type) {
	funcValue = reflect.ValueOf(outFunc)
	funcType := funcValue.Type()
	if funcType.Kind() != reflect.Func ||
		funcType.IsVariadic() ||
		funcType.NumIn() != 3 ||
		funcType.In(0) != resourceContextPtrType ||
		funcType.In(1).Kind() != reflect.Struct ||
		funcType.In(2).Kind() != reflect.Struct ||
		funcType.NumOut() != 2 ||
		funcType.Out(0).Kind() != reflect.Struct ||
		funcType.Out(1) != errorType {
		panic(fmt.Sprintf("outFunc must have signature like '%s', got '%T'", OutFuncPattern, outFunc))
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
	version, err := call(funcValue, reflect.ValueOf(&resourceContext), sourceValue, paramsValue)
	if err != nil {
		return err
	}

	return json.NewEncoder(respWriter).Encode(ResourceResponse{
		Version:  version,
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

func call(funcValue reflect.Value, in ...reflect.Value) (result interface{}, err error) {
	resultValues := funcValue.Call(in)
	result = resultValues[0].Interface()
	if !resultValues[1].IsNil() {
		err = resultValues[1].Interface().(error)
	}
	return
}
