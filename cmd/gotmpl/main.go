package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joshuagrisham-karolinska/gotmpl"
	"github.com/joshuagrisham-karolinska/gotmpl/template"
)

func main() {

	// Set up and parse options
	usage := `Render a Go text template using the given data file.
Usage:
  gotmpl (--template <path> --data <path>)
  gotmpl --help
  gotmpl --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -t --template <path>  Template file path.
  -d --data <path>      Data file path.`

	opts, _ := docopt.ParseArgs(usage, os.Args[1:], gotmpl.Version)
	tmplPath, _ := opts.String("--template")
	dataPath, _ := opts.String("--data")

	// Read template from file system
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Read data from file system
	dataBytes, err := os.ReadFile(dataPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Unmarshal data JSON string to a map[string]interface{}
	// TODO: Support YAML in addition to JSON?
	data := make(map[string]interface{})
	err = json.Unmarshal(dataBytes, &data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Render template using data and write the result to os.Stdout
	err = template.Render(string(tmplBytes), data, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
