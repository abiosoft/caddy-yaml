# caddy-yaml

[![Go](https://github.com/abiosoft/caddy-yaml/workflows/Go/badge.svg)](https://github.com/abiosoft/caddy-yaml/actions)


Alternative Caddy YAML config adapter with extra features.

## Install

Install with [xcaddy](https://github.com/caddyserver/xcaddy).

```
xcaddy build \
    --with github.com/abiosoft/caddy-yaml
```
## Usage

Specify with the `--adapter` flag for `caddy run`.
```
caddy run --config /path/to/yaml/config.yaml --adapter yaml
```

## Comparison with [iamd3vil/caddy_yaml_adapter](https://github.com/iamd3vil/caddy_yaml_adapter)

This project does a few extra things.

* Conditional configurations with Templates.

  ```yaml
  #{if ne $ENVIRONMENT "production"}
  logging:
    logs:
      default: { level: DEBUG }
  #{end}
  ```

* Top level configs prefixed with `_` are discarded.

  This makes it possible to leverage features like YAML anchors. Otherwise, Caddy would error for unknown fields.

  ```yaml
  ...
  _domain: &domain mysite.example.com
  ...
  host: [ *domain ]
  ...
  logger_names: 
    - *domain: customlog
  ...
  ```

* Config-time environment variables

  Without the Caddyfile, Caddy's native configuration limits to [runtime environment variables](https://caddyserver.com/docs/caddyfile/concepts#environment-variables).
  There are use cases for knowing the environment variables at configuration time. e.g. troubleshooting purposes.

  ```yaml
  listen: "#{ $PORT }"
  ```

If the above features are not needed, the behaviour is identical to [iamd3vil/caddy_yaml_adapter](https://github.com/iamd3vil/caddy_yaml_adapter).


**Note** that you can not have both adapters built with Caddy, they are incompatible. They both register as `yaml` config adapter and at most one config adapter is allowed per config format.


## Templating

Anything supported by [Go templates](https://pkg.go.dev/text/template) can be used, as well as any [Sprig](https://masterminds.github.io/sprig) function.

### Delimeters

Delimeters are `#{` and `}`. e.g. `#{.text}`. This ensures the YAML config file with templates included remains a valid YAML file and can still be validated by the schema.

### Values

Top level configs prefixed with `_` can be reused anywhere else in
the YAML config.

```yaml
_hello: Hello from YAML template
_nest:
  value: nesting
```

Referencing in a route handler.

```yaml
...
handle:
  - handler: static_response
    body: "#{ ._hello } with #{ ._nest.value }"
```

### Environment Variables

Environment variables can be referenced in a template by prefixing with `$` sign.

```yaml
listen:
    - "#{ $PORT }"
...
handler: file_server
root: "#{ $PROJECT_DIRECTORY }/public"
```

Caddy supports runtime environment variables via [`{env.*}` placeholders](https://caddyserver.com/docs/caddyfile/concepts#environment-variables).

### Example Config 

Check the [test YAML configuration file](https://github.com/abiosoft/caddy-yaml/blob/master/testdata/test.caddy.yaml).

## License

Apache 2
