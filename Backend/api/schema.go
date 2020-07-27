package api

import (
	"bytes"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var cache = make(map[string]*gojsonschema.Schema)

func validateSchema(request *http.Request, schemaText string) (error, io.Reader) {
	var err error

	schema, ok := cache[schemaText]
	if !ok {
		loader := gojsonschema.NewStringLoader(schemaText)
		sl := gojsonschema.NewSchemaLoader()
		schema, err = sl.Compile(loader)
		if err != nil {
			return err, nil
		}

		cache[schemaText] = schema
	}

	buf, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return err, nil
	}

	document := gojsonschema.NewBytesLoader(buf)
	result, err := schema.Validate(document)
	if err != nil {
		return err, nil
	}

	if !result.Valid() {
		return &ValidationError{result}, nil
	}

	buffer := bytes.NewBuffer(buf)
	return nil, buffer
}

type ValidationError struct {
	Result *gojsonschema.Result
}

func (error ValidationError) Error() string {
	errors := make([]string, len(error.Result.Errors()))

	for _, error := range error.Result.Errors() {
		errors = append(errors, error.String())
	}
	return strings.Join(errors, "\n")
}
