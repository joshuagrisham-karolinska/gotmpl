package template

import (
	"io"
	"text/template"
)

// Hard-coded "name" for the temporary Template instance that will be created
const templateName = "gotmpl"

// Creates a temporary instance of a Text Template based on a string-representation of the desired template,
// executes the template using the given data interface{}, and writes the result to the given Writer.
// - sprig v3 functions are available
// - missing keys will result in an error
func Render(tmpl string, data interface{}, w io.Writer) error {
	return template.Must(template.New(templateName).Option("missingkey=error").Funcs(funcMap()).Parse(tmpl)).Execute(w, data)
}
