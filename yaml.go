package caddyyaml

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func init() {
	caddyconfig.RegisterAdapter("yaml", Adapter{})
}

// Adapter adapts YAML to Caddy JSON.
type Adapter struct{}

// Adapt converts the YAML config in body to Caddy JSON.
func (a Adapter) Adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	return adapt(body, options)
}
