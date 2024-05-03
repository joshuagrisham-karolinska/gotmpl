//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"encoding/json"
	"syscall/js"

	"github.com/joshuagrisham-karolinska/gotmpl/template"
)

func Render(this js.Value, args []js.Value) interface{} {

	result := make(map[string]interface{})

	// TODO: In case text/template execution fails it will panic
	// For now just return gracefully instead of letting the program terminate
	// Better would be to potentially use some kind of reflection to get the actual error
	// message from the panic and then return that error instead
	defer func() {
		if err := recover(); err != nil {
			return
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
