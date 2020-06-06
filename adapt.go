package caddyyaml

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	// extract variables
	vars, err := varsFromBody(body)
	if err != nil {
		return nil, nil, err
	}

	// apply template
	tmp, err := applyTemplate(body, vars)
	if err != nil {
		return nil, nil, err
	}

	// convert to YAML
	result, err := yamlToJSON(tmp)
	return result, nil, err
}
