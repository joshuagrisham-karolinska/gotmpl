/*
Tests from Helm are copyright The Helm Authors.
*/

package template

import (
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestFuncs(t *testing.T) {

	tests := []struct {
		tpl, expect string
		vars        interface{}
	}{{
		tpl:    `{{ toUTCDateTime .str "Europe/Stockholm" }}`,
		expect: `2023-07-08T16:09:43.345Z`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43.345+02:00"},
	}, {
		tpl:    `{{ toUTCDateTime .str }}`,
		expect: `2023-07-08T18:09:43.123Z`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43.123Z"},
	}, {
		tpl:    `{{ toUTCDateTime .str }}`,
		expect: `2023-07-08T16:09:43.123Z`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43.123+02:00"},
	}, {
		tpl:    `{{ toUTCDateTime .str "America/New_York" }}`,
		expect: `2023-07-08T16:09:43.123Z`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43.123+02:00"},
	}, {
		tpl:    `{{ toUTCDateTime .str "America/New_York" }}`,
		expect: `2023-07-08T22:09:43.123Z`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43.123"},
	}, {
		tpl:    `{{ toLocalDateTime .str }}`,
		expect: `2023-07-08T20:09:43.123+02:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43.123Z"},
	}, {
		tpl:    `{{ toLocalDateTime .str }}`,
		expect: `2023-07-08T18:09:43.000+02:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str }}`,
		expect: `2023-07-08T18:09:43.000+02:00`,
		vars:   map[string]interface{}{"str": "2023-07-08T18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" }}`,
		expect: `2023-07-09T00:09:43.000+02:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "Australia/Sydney" }}`,
		expect: `2023-07-09T08:09:43.000+10:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "-0800" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "-08" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "-8" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "UTC-8" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "UTC-0800" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43"},
	}, {
		tpl:    `{{ toLocalDateTime .str "America/New_York" "UTC -08:00" }}`,
		expect: `2023-07-08T14:09:43.000-08:00`,
		vars:   map[string]interface{}{"str": "2023-07-08 18:09:43 EDT"},
	}, {
		tpl:    `{{ toUTCDateTime .str }}`,
		expect: `2023-07-08T00:00:00.000Z`,
		vars:   map[string]interface{}{"str": "2023-07-08"},
	}}

	for _, tt := range tests {
		var b strings.Builder
		err := template.Must(template.New("test").Funcs(funcMap()).Parse(tt.tpl)).Execute(&b, tt.vars)
		assert.NoError(t, err)
		assert.Equal(t, tt.expect, b.String(), tt.tpl)
	}
}

// Below tests were taken from https://github.com/helm/helm/blob/588041f6a55a8f23113f3c080d1733e65176fa0a/pkg/engine/funcs_test.go
// but removed the tests for toJson / fromJson as these are already included in Sprig
func TestFuncsFromHelm(t *testing.T) {
	//TODO write tests for failure cases
	tests := []struct {
		tpl, expect string
		vars        interface{}
	}{{
		tpl:    `{{ toYaml . }}`,
		expect: `foo: bar`,
		vars:   map[string]interface{}{"foo": "bar"},
	}, {
		tpl:    `{{ toToml . }}`,
		expect: "foo = \"bar\"\n",
		vars:   map[string]interface{}{"foo": "bar"},
	}, {
		tpl:    `{{ fromYaml . }}`,
		expect: "map[hello:world]",
		vars:   `hello: world`,
	}, {
		tpl:    `{{ fromYamlArray . }}`,
		expect: "[one 2 map[name:helm]]",
		vars:   "- one\n- 2\n- name: helm\n",
	}, {
		tpl:    `{{ fromYamlArray . }}`,
		expect: "[one 2 map[name:helm]]",
		vars:   `["one", 2, { "name": "helm" }]`,
	}, {
		// Regression for https://github.com/helm/helm/issues/2271
		tpl:    `{{ toToml . }}`,
		expect: "[mast]\n  sail = \"white\"\n",
		vars:   map[string]map[string]string{"mast": {"sail": "white"}},
	}, {
		tpl:    `{{ fromYaml . }}`,
		expect: "map[Error:error unmarshaling JSON: while decoding JSON: json: cannot unmarshal array into Go value of type map[string]interface {}]",
		vars:   "- one\n- two\n",
	}, {
		tpl:    `{{ fromJsonArray . }}`,
		expect: `[one 2 map[name:helm]]`,
		vars:   `["one", 2, { "name": "helm" }]`,
	}, {
		tpl:    `{{ fromJsonArray . }}`,
		expect: `[json: cannot unmarshal object into Go value of type []interface {}]`,
		vars:   `{"hello": "world"}`,
	}, {
		tpl:    `{{ merge .dict (fromYaml .yaml) }}`,
		expect: `map[a:map[b:c]]`,
		vars:   map[string]interface{}{"dict": map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, "yaml": `{"a":{"b":"d"}}`},
	}, {
		tpl:    `{{ merge (fromYaml .yaml) .dict }}`,
		expect: `map[a:map[b:d]]`,
		vars:   map[string]interface{}{"dict": map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, "yaml": `{"a":{"b":"d"}}`},
	}, {
		tpl:    `{{ fromYaml . }}`,
		expect: `map[Error:error unmarshaling JSON: while decoding JSON: json: cannot unmarshal array into Go value of type map[string]interface {}]`,
		vars:   `["one", "two"]`,
	}, {
		tpl:    `{{ fromYamlArray . }}`,
		expect: `[error unmarshaling JSON: while decoding JSON: json: cannot unmarshal object into Go value of type []interface {}]`,
		vars:   `hello: world`,
	}}

	for _, tt := range tests {
		var b strings.Builder
		err := template.Must(template.New("test").Funcs(funcMap()).Parse(tt.tpl)).Execute(&b, tt.vars)
		assert.NoError(t, err)
		assert.Equal(t, tt.expect, b.String(), tt.tpl)
	}
}

// This test to check a function provided by sprig is due to a change in a
// dependency of sprig. mergo in v0.3.9 changed the way it merges and only does
// public fields (i.e. those starting with a capital letter). This test, from
// sprig, fails in the new version. This is a behavior change for mergo that
// impacts sprig and Helm users. This test will help us to not update to a
// version of mergo (even accidentally) that causes a breaking change. See
// sprig changelog and notes for more details.
// Note, Go modules assume semver is never broken. So, there is no way to tell
// the tooling to not update to a minor or patch version. `go install` could
// be used to accidentally update mergo. This test and message should catch
// the problem and explain why it's happening.
func TestMergeFromHelm(t *testing.T) {
	dict := map[string]interface{}{
		"src2": map[string]interface{}{
			"h": 10,
			"i": "i",
			"j": "j",
		},
		"src1": map[string]interface{}{
			"a": 1,
			"b": 2,
			"d": map[string]interface{}{
				"e": "four",
			},
			"g": []int{6, 7},
			"i": "aye",
			"j": "jay",
			"k": map[string]interface{}{
				"l": false,
			},
		},
		"dst": map[string]interface{}{
			"a": "one",
			"c": 3,
			"d": map[string]interface{}{
				"f": 5,
			},
			"g": []int{8, 9},
			"i": "eye",
			"k": map[string]interface{}{
				"l": true,
			},
		},
	}
	tpl := `{{merge .dst .src1 .src2}}`
	var b strings.Builder
	err := template.Must(template.New("test").Funcs(funcMap()).Parse(tpl)).Execute(&b, dict)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"a": "one", // key overridden
		"b": 2,     // merged from src1
		"c": 3,     // merged from dst
		"d": map[string]interface{}{ // deep merge
			"e": "four",
			"f": 5,
		},
		"g": []int{8, 9}, // overridden - arrays are not merged
		"h": 10,          // merged from src2
		"i": "eye",       // overridden twice
		"j": "jay",       // overridden and merged
		"k": map[string]interface{}{
			"l": true, // overridden
		},
	}
	assert.Equal(t, expected, dict["dst"])
}
