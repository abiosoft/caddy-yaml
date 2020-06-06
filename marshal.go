package caddyyaml

import (
	"encoding/json"
	"strings"

	"github.com/ghodss/yaml"
)

func yamlToJSON(b []byte) ([]byte, error) {
	var tmp map[string]interface{}

	if err := yaml.Unmarshal(b, &tmp); err != nil {
		return nil, err
	}

	// discard all entries with _ prefix
	for key := range tmp {
		if strings.HasPrefix(key, "_") {
			// this is safe to do
			// https://stackoverflow.com/a/23230406/524060
			delete(tmp, key)
		}
	}

	return json.Marshal(tmp)
}

func varsFromBody(b []byte) (map[string]interface{}, error) {
	var vars map[string]interface{}

	varsBytes, err := extractVariables(b)
	if err != nil {
		return nil, err
	}

	varsBytes, err = applyTemplate(varsBytes, nil)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(varsBytes, &vars); err != nil {
		return nil, err
	}

	return vars, nil

}
