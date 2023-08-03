package caddyyaml

import (
	"github.com/caddyserver/caddy/v2/caddyconfig"
)

func init() {
	caddyconfig.RegisterAdapter("yaml", Adapter{})
}

// Adapter adapts YAML to Caddy JSON.
type Adapter struct{}

// envOptionName is the name of the option to set an environment array in the options.
// This is mainly intended as an internal option to aid in testing, if this is not set then
// the value of `os.Environ()` is used.
const envOptionName = "yaml.Env"

// Adapt converts the YAML config in body to Caddy JSON.
func (a Adapter) Adapt(body []byte, options map[string]interface{}) ([]byte, []caddyconfig.Warning, error) {
	return adapt(body, options)
}
