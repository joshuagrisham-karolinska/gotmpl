package template

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Masterminds/sprig/v3"
	"sigs.k8s.io/yaml"
)

// funcMap returns a mapping of all of the functions that Engine has.
func funcMap() template.FuncMap {
	// use Sprig's TxtFuncMap as a base
	f := sprig.TxtFuncMap()

	// remove environment variable stuff -- these should not be used for our case
	delete(f, "env")
	delete(f, "expandenv")

	// Add some extra functionality
	extra := template.FuncMap{
		"toUTCDateTime":   toUTCDateTime,
		"toLocalDateTime": toLocalDateTime,

		"toToml":        toTOML,
		"toYaml":        toYAML,
		"fromYaml":      fromYAML,
		"fromYamlArray": fromYAMLArray,
		"fromJsonArray": fromJSONArray,
	}

	// add each entry in `extra` to `f`
	for k, v := range extra {
		f[k] = v
	}

	return f
}

//*** Date functions ***//

// toUTCDateTime converts many recognized datetime string formats to an ISO8601-formatted string in UTC time zone
//
// An optional second string parameter can be provided as a time zone location name or offset to interpret
// the incoming datetime str with, but will be used only in case the provided datetime str value does not
// already specify the time offset information in one of the recognized formats
func toUTCDateTime(str string, locationIn ...string) string {
	intzstr := "UTC"
	if len(locationIn) > 0 && strings.TrimSpace(locationIn[0]) != "" {
		intzstr = locationIn[0]
	}
	return toLocalDateTime(str, intzstr, "UTC")
}

// toLocalDateTime converts many recognized datetime string formats to an ISO8601-formatted string in a local time zone
//
// The optional second string parameter can be provided as a time zone location name or offset to interpret
// the incoming datetime str with, but will be used only in case the provided datetime str value does not
// already specify the time offset information in one of the recognized formats. If omitted, a system default will be used.
//
// The optional third string parameter can be provided as the desired target time zone to convert the input
// datetime string to. If omitted, a system default will be used.
func toLocalDateTime(str string, locationInOut ...string) string {

	// "Local" defaults to Europe/Stockholm TODO: set via env or config??
	intz, _ := time.LoadLocation("Europe/Stockholm")
	outtz, _ := time.LoadLocation("Europe/Stockholm")
	if len(locationInOut) > 0 && strings.TrimSpace(locationInOut[0]) != "" {
		parsedtz, err := parseTzOffset(locationInOut[0])
		if err != nil {
			return ""
		}
		intz = parsedtz
	}
	if len(locationInOut) > 1 && strings.TrimSpace(locationInOut[1]) != "" {
		parsedtz, err := parseTzOffset(locationInOut[1])
		if err != nil {
			return ""
		}
		outtz = parsedtz
	}

	// ordered list of time formats to attempt to match
	// the function will return the value which first successfully parses, otherwise "zero time" will be returned
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.000Z07:00",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05.000",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05.000 MST",
		"2006-01-02 15:04:05 MST",
		time.DateTime,
		time.DateOnly,
		time.Layout,
	}

	var result time.Time
	for _, format := range formats {
		parsed, err := time.ParseInLocation(format, str, intz)
		if err == nil {
			result = parsed
		}
	}

	return result.In(outtz).Format("2006-01-02T15:04:05.000Z07:00")
}

// parseTzOffset is a helper function to parse a Location name or offset (e.g. "UTC+2", "-0700", etc) and return as a time.Location
func parseTzOffset(str string) (*time.Location, error) {

	// first see if the given str is a valid Location name; if so, return it
	location, err := time.LoadLocation(str)
	if err == nil {
		return location, nil
	}

	// otherwise, try to parse the value as an offset integer and return it as new Location using time.FixedZone()

	// get the sign and strip the prefix from the offset amount
	var sign int
	var offsetstr string
	prefixes := []struct {
		prefix string
		sign   int
	}{
		{prefix: "UTC-", sign: -1},
		{prefix: "UTC+", sign: 1},
		{prefix: "UTC -", sign: -1},
		{prefix: "UTC +", sign: 1},
		{prefix: "-", sign: -1},
		{prefix: "+", sign: 1},
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(str, prefix.prefix) {
			sign = prefix.sign
			offsetstr = strings.TrimPrefix(str, prefix.prefix)
			break
		}
	}
	if offsetstr == "" {
		return nil, &time.ParseError{Message: "could not parse offset string"}
	}

	// now try to parse the offset amount as a time (hours and minutes)
	var offsettime time.Time
	offsetformats := []string{
		"1504",
		"15:04",
		"15",
	}
	for _, offsetformat := range offsetformats {
		parsedtime, err := time.Parse(offsetformat, offsetstr)
		if err == nil {
			offsettime = parsedtime
		}
	}
	if offsettime.IsZero() {
		return nil, &time.ParseError{Message: "could not parse offset time"}
	}

	// and finally compute the offset amount as +/-(hours + minutes + seconds)
	offsetint := sign * ((offsettime.Hour() * 60 * 60) + (offsettime.Minute() * 60) + offsettime.Second())
	return time.FixedZone(str, offsetint), nil

}

// Copy some functions from Helm just to add feature parity
// see: https://github.com/helm/helm/blob/588041f6a55a8f23113f3c080d1733e65176fa0a/pkg/engine/funcs.go

// toYAML takes an interface, marshals it to yaml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toYAML(v interface{}) string {
	data, err := yaml.Marshal(v)
	if err != nil {
		// Swallow errors inside of a template.
		return ""
	}
	return strings.TrimSuffix(string(data), "\n")
}

// fromYAML converts a YAML document into a map[string]interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string into
// m["Error"] in the returned map.
func fromYAML(str string) map[string]interface{} {
	m := map[string]interface{}{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		m["Error"] = err.Error()
	}
	return m
}

// fromYAMLArray converts a YAML array into a []interface{}.
//
// This is not a general-purpose YAML parser, and will not parse all valid
// YAML documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string as
// the first and only item in the returned array.
func fromYAMLArray(str string) []interface{} {
	a := []interface{}{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}

// toTOML takes an interface, marshals it to toml, and returns a string. It will
// always return a string, even on marshal error (empty string).
//
// This is designed to be called from a template.
func toTOML(v interface{}) string {
	b := bytes.NewBuffer(nil)
	e := toml.NewEncoder(b)
	err := e.Encode(v)
	if err != nil {
		return err.Error()
	}
	return b.String()
}

// Note: toJson and fromJson are already included in Sprig

// fromJSONArray converts a JSON array into a []interface{}.
//
// This is not a general-purpose JSON parser, and will not parse all valid
// JSON documents. Additionally, because its intended use is within templates
// it tolerates errors. It will insert the returned error message string as
// the first and only item in the returned array.
func fromJSONArray(str string) []interface{} {
	a := []interface{}{}

	if err := json.Unmarshal([]byte(str), &a); err != nil {
		a = []interface{}{err.Error()}
	}
	return a
}
