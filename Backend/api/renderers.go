package api

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"io"
	"net/http"
)

func RenderJson(writer io.Writer, data interface{}) {
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		fmt.Println("failed to write json")
	}
}

func renderError(writer http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *ValidationError:
		writer.WriteHeader(400)

		errors := make([]interface{}, 0)

		for _, result := range e.Result.Errors() {
			switch r := result.(type) {
			case *gojsonschema.RequiredError:
				errors = append(errors, map[string]interface{}{"code": "missing_field", "field": r.Details()["property"]})
			default:
				errors = append(errors, map[string]interface{}{"code": r.Type()})
			}

		}
		RenderJson(writer, map[string]interface{}{
			"code":   "invalid_request",
			"errors": errors,
		})

	default:
		fmt.Println(err)
		writer.WriteHeader(500) // todo: use 400 for jsonschema errors
		RenderJson(writer, map[string]interface{}{"code": "server_error"})
	}
}
