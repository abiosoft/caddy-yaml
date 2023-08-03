package caddyyaml

import (
	"errors"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	filename, ok := options["filename"].(string)
	if !ok {
		return nil, nil, errors.New("missing filename option")
	}

	wc := newWarningsCollector(filename)

	// extract variables
	vars, err := varsFromBody(body)
	if err != nil {
		return nil, wc.warnings, err
	}

	// apply template
	tmp, err := applyTemplate(body, vars)
	if err != nil {
		return nil, wc.warnings, err
	}

	// convert to YAML
	result, err := yamlToJSON(tmp)
	return result, wc.warnings, err
}
