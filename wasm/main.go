//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/clbanning/mxj/v2"
	"github.com/joshuagrisham-karolinska/gotmpl/template"
	"sigs.k8s.io/yaml"
)

func Render(this js.Value, args []js.Value) (value any) {

	result := make(map[string]interface{})

	// If text/template execution fails it will panic
	// In this case we will capture the panic and return it as the "errorTmpl" value
	defer func() {
		result := make(map[string]interface{})
		if err := recover(); err != nil {
			result["errorTmpl"] = fmt.Sprintf("%v", err)
			value = result
		}
	}()

	tmpl := args[0].String()
	data := strings.TrimSpace(args[1].String())

	result["data"] = data
	result["tmpl"] = tmpl

	dataMap := make(map[string]interface{})
	var err error

	// TODO: Maybe better to add another argument for setting the format (JSON vs YAML vs XML)
	// for now we will write some logic to try and "guess" JSON vs YAML vs XML by looking at the first character of the data

	switch {

	// unmarshall to map[string]interface requires an object at the top level even though valid JSON can start with an array or a single element value
	// so here we will only try to detect if the data starts with "{" and assume it will be JSON
	case data[:1] == "{":
		err = json.Unmarshal([]byte(data), &dataMap)

	// beginning with "<" is assumed to be XML
	case data[:1] == "<":
		dataMap, err = mxj.NewMapXml([]byte(data))

	// otherwise we can just assume it is YAML, I guess ? (since in YAML you can quote key names and stuff..)
	default:
		err = yaml.Unmarshal([]byte(data), &dataMap)

	}

	result["dataMap"] = dataMap
	if err != nil {
		result["errorData"] = err.Error()
		return result
	}

	// create a buffer for template.Render to write its results to
	buf := new(bytes.Buffer)

	// Render template using dataMap and write the result to the buffer
	err = template.Render(tmpl, dataMap, buf)
	if err != nil {
		result["errorTmpl"] = err.Error()
		return result
	}

	result["output"] = buf.String()
	return result

}

func main() {
	//c := make(chan bool)
	// Register JavaScript reference to Go function
	js.Global().Set("render", js.FuncOf(Render))
	// Stay loaded indefinitely
	select {}
	//<-c
}
