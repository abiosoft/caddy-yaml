package caddyyaml

import (
	"fmt"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/ghodss/yaml"
)

func adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	// split config from body
	body, configBody, err := extractConfigs(body, templateValuesKey)
	if err != nil {
		return configBody, nil, err
	}

	// apply template for env vars.
	configBody, err = applyTemplate(configBody, nil)
	if err != nil {
		return configBody, nil, err
	}
	// marshal config
	config := map[string]interface{}{}
	if err := yaml.Unmarshal(configBody, &config); err != nil {
		return nil, nil, err
	}

	// if no config, convert to yaml as is.
	values, ok := config[templateValuesKey]
	if !ok {
		result, err := yaml.YAMLToJSON(body)
		return result, nil, err
	}

	// validate template values
	if _, ok := values.(map[string]interface{}); !ok {
		return nil, nil, fmt.Errorf("%s must be a map", templateValuesKey)
	}

	// apply template
	tmp, err := applyTemplate(body, values.(map[string]interface{}))
	if err != nil {
		return nil, nil, err
	}

	// convert to YAML
	result, err := yaml.YAMLToJSON(tmp)
	return result, nil, err
}
