package caddyyaml

import (
	"errors"
	"os"

	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	filename, ok := options["filename"].(string)
	if !ok {
		return nil, nil, errors.New("missing filename option")
	}

	env, ok := options[envOptionName].([]string)
	if !ok {
		env = os.Environ()
	}

	wc := newWarningsCollector(filename)

	// extract variables
	vars, err := varsFromBody(body, env)
	if err != nil {
		return nil, wc.warnings, err
	}

	// apply template
	tmp, err := applyTemplate(body, vars, env, wc)
	if err != nil {
		return nil, wc.warnings, err
	}

	// convert to YAML
	result, err := yamlToJSON(tmp)
	return result, wc.warnings, err
}
