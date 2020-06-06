# caddy-yaml
Alternative Caddy YAML config adapter with template support.

## Install

```
xcaddy build \
    --with github.com/abiosoft/caddy-yaml
```

## Usage

Anything supported by [Go templates](https://pkg.go.dev/text/template) can be used, as well as any [sprig](https://masterminds.github.io/sprig) functions.

### Delimeters

Delimeters are `#{` and `}`. e.g. `#{.text}`. This ensures the YAML config file with templates included remains a valid YAML file and can still be validated by the schema.

### Values

Anything defined in top level `template_values` can be reused anywhere else in
the YAML config.

```yaml
template_values:
  hello: Hello from YAML template
  nest:
    value: nesting
```
Referencing in a route handler.

```yaml
...
handle:
  - handler: static_response
    body: "#{ .hello } with #{ .nest.value }"
```

### Environment Variables

Simply reference the environment variable in the template.

```yaml
listen:
    - "#{ $PORT }"
...
handler: file_server
root: "#{ $PROJECT }/public"
```

Caddy also supports runtime environment variables via [`{env.*}` placeholders](https://caddyserver.com/docs/caddyfile/concepts#environment-variables).

### Example

Check the [test file](https://github.com/abiosoft/caddy-yaml/blob/master/testdata/test.caddy.yaml).

## Comparison with iamd3vil/caddy_yaml_adapter

If you are not using templates, the behaviour is identical to [iamd3vil/caddy_yaml_adapter](https://github.com/iamd3vil/caddy_yaml_adapter).

**Note** that you can not have both adapters built with Caddy, they are incompatible. They both register as `yaml` config adapter and there can be at most one config adapter per config format.

## License

Apache 2