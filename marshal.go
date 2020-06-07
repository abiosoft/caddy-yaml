package caddyyaml

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

func yamlToJSON(b []byte) ([]byte, error) {
	var tmp map[string]interface{}

	if err := yaml.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	// discard all entries with x- prefix
	for key := range tmp {
		if strings.HasPrefix(key, "x-") {
			// this is safe to do
			// https://stackoverflow.com/a/23230406/524060
			delete(tmp, key)
		}
	}

	return json.Marshal(tmp)
}

func varsFromBody(b []byte) (map[string]interface{}, error) {
	var tmp map[string]interface{}
	var vars map[string]interface{}

	varsBytes, err := extractVariables(b)
	if err != nil {
		return nil, err
	}

	varsBytes, err = applyTemplate(varsBytes, nil)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(varsBytes, &tmp); err != nil {
		return nil, err
	}

	vars = make(map[string]interface{})

	// go template prohibits hyphen `-` in field names.
	for xkey, val := range tmp {
		key := xkey[2:] // discard x- prefix
		if strings.Index(key, "-") > 0 {
			return nil, fmt.Errorf("template: apart from 'x-' prefix, '-' cannot be used in extension field name for %s", xkey)
		}
		vars[key] = val
	}

	return vars, nil

}
