package caddyyaml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/ghodss/yaml"
)

func init() {
	caddyconfig.RegisterAdapter("yaml", Adapter{})
}

const templateObjKey = "template_values"

// Adapter adapts YAML to Caddy JSON.
type Adapter struct{}

// Adapt converts the YAML config in body to Caddy JSON.
func (a Adapter) Adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	// split config from body
	body, configBody, err := extractConfigs(body, templateObjKey)
	if err != nil {
		return configBody, nil, err
	}

	// marshal config
	config := map[string]interface{}{}
	if err := yaml.Unmarshal(configBody, &config); err != nil {
		return nil, nil, err
	}

	// if no config, convert to yaml as is.
	values, ok := config[templateObjKey]
	if !ok {
		result, err := yaml.YAMLToJSON(body)
		return result, nil, err
	}

	// validate values
	if _, ok := values.(map[string]interface{}); !ok {
		return nil, nil, fmt.Errorf("%s must be a map", templateObjKey)
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

func applyTemplate(body []byte, values map[string]interface{}) ([]byte, error) {
	tpl, err := template.New("yaml").Delims("#{", "}").Parse(string(body))
	if err != nil {
		return body, err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, values); err != nil {
		return body, err
	}
	return out.Bytes(), nil
}

func stripFields(body []byte, fields ...string) ([]byte, error) {
	var obj map[string]interface{}
	if err := yaml.Unmarshal(body, &obj); err != nil {
		return body, err
	}
	for _, field := range fields {
		delete(obj, field)
	}

	return json.Marshal(obj)
}
