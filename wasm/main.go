//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/joshuagrisham-karolinska/gotmpl/template"
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
	data := args[1].String()

	result["data"] = data
	result["tmpl"] = tmpl

	// Unmarshal data JSON string to a map[string]interface{}
	// TODO: Support YAML in addition to JSON?
	dataMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &dataMap)
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
