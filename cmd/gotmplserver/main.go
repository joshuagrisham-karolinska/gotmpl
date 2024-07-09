package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/clbanning/mxj/v2"
	"github.com/docopt/docopt-go"
	"github.com/joshuagrisham-karolinska/gotmpl"
	"github.com/joshuagrisham-karolinska/gotmpl/template"
	"sigs.k8s.io/yaml"
)

func main() {

	// Set up and parse options
	usage := `Start gotmpl HTTP server.
Usage:
  gotmplserver
  gotmplserver [--port <port> --path <path>]
  gotmplserver --help | --version

Options:
  -h --help         Show this screen.
  -v --version      Show version.
  -p --port <port>  HTTP port number [default: 10000].
  --path <path>     HTTP path [default: /gotmpl].`

	opts, _ := docopt.ParseArgs(usage, os.Args[1:], gotmpl.Version)
	port, _ := opts.String("--port")
	path, _ := opts.String("--path")

	log.Printf("Starting gotmpl Server; listening on http://0.0.0.0:%s%s\n", port, path)

	// Set up HTTP handler functions and start the server
	http.HandleFunc(path, handlePath)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func handlePath(w http.ResponseWriter, r *http.Request) {

	// Only allow POST method
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", "POST")
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}

	// If text/template execution fails it will panic
	// In this case we will capture the panic and return an error with reason TemplateError
	defer func() {
		if err := recover(); err != nil {
			writeHttpBadRequest(w, "TemplateError", fmt.Sprintf("%v", err))
			return
		}
	}()

	// Get template and data from form values
	// Note: per https://pkg.go.dev/net/http#Request.FormValue
	//  using FormValue supports the client to set these values in any of:
	//    - application/x-www-form-urlencoded
	//    - query parameters
	//    - multipart/form-data
	tmpl := r.FormValue("template")
	dataString := strings.TrimSpace(r.FormValue("data"))

	data := make(map[string]interface{})
	var err error

	// TODO: Maybe better to support setting and reading Content-Type per part of multipart/form-data ?
	// for now we will write some logic to try and "guess" JSON vs YAML vs XML by looking at the first character of the data

	switch {

	// unmarshall to map[string]interface requires an object at the top level even though valid JSON can start with an array or a single element value
	// so here we will only try to detect if the data starts with "{" and assume it will be JSON
	case dataString[:1] == "{":
		err = json.Unmarshal([]byte(dataString), &data)

	// beginning with "<" is assumed to be XML
	case dataString[:1] == "<":
		data, err = mxj.NewMapXml([]byte(dataString))

	// otherwise we can just assume it is YAML, I guess ? (since in YAML you can quote key names and stuff..)
	default:
		err = yaml.Unmarshal([]byte(dataString), &data)

	}

	if err != nil {
		writeHttpBadRequest(w, "DataUnmarshallingError", err.Error())
		return
	}

	// Render template using data and write the result to the ResponseWriter
	err = template.Render(tmpl, data, w)
	if err != nil {
		writeHttpBadRequest(w, "TemplateRenderingError", err.Error())
		return
	}

}

type HttpError struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type HttpErrorResponse struct {
	Error HttpError `json:"error"`
}

func writeHttpBadRequest(w http.ResponseWriter, reason string, message string) {
	response := HttpErrorResponse{HttpError{reason, message}}
	responseBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(responseBytes)
}
