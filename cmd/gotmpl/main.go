package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/clbanning/mxj/v2"
	"github.com/docopt/docopt-go"
	"github.com/joshuagrisham-karolinska/gotmpl"
	"github.com/joshuagrisham-karolinska/gotmpl/template"
	"sigs.k8s.io/yaml"
)

func main() {

	// Set up and parse options
	usage := `Render a Go text template using the given data file.
Usage:
  gotmpl (--template <path> --data <path>)
  gotmpl --help | --version

Options:
  -h --help             Show this screen.
  -v --version          Show version.
  -t --template <path>  Template file path.
  -d --data <path>      Data file path (supports JSON, YAML, and XML).`

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

	data := make(map[string]interface{})
	switch filepath.Ext(dataPath) {
	case ".json":
		err = json.Unmarshal(dataBytes, &data)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(dataBytes, &data)
	case ".xml":
		// use mxj instead of encoding/xml since we want to use generic map[string]interface
		data, err = mxj.NewMapXml(dataBytes)
	default:
		fmt.Printf("unsupported data file extension '%s'\n", filepath.Ext(dataPath))
		os.Exit(1)
	}
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
