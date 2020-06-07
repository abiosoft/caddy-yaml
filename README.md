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

## Comparison with existing YAML adapter

This project does a few extra things.

* Conditional configurations with Templates.

  ```yaml
  #{if ne $ENVIRONMENT "production"}
  logging:
    logs:
      default: { level: DEBUG }
  #{end}
  ```

* Extension Fields. A convension used in [Docker Compose](https://docs.docker.com/compose/compose-file/#extension-fields). 

  Top level keys prefixed with `x-` are discarded. This makes it easier to leverage YAML anchors and aliases, while avoiding Caddy errors due to unknown fields.

  ```yaml
  # anchor declaration
  x-domain: &domain mysite.example.com
  # reuse alias
  host: [ *domain ]
  # reuse alias
  logger_names: 
    - *domain: customlog
  ```

* Config-time environment variables

  Without the Caddyfile, Caddy's native configuration limits to runtime environment variables.
  There are use cases for knowing the environment variables at configuration time. e.g. troubleshooting purposes.

  ```yaml
  listen: "#{ $PORT }"
  ```

If the above features are not needed or utilised, the behaviour is identical to [iamd3vil/caddy_yaml_adapter](https://github.com/iamd3vil/caddy_yaml_adapter).


_**Note** that both adapters cannot be built with Caddy, they are incompatible. They both register as `yaml` config adapter and at most one config adapter is allowed per config format_.


## Templating

Anything supported by [Go templates](https://pkg.go.dev/text/template) can be used, as well as any [Sprig](https://masterminds.github.io/sprig) function.

### Delimeters

Delimeters are `#{` and `}`. e.g. `#{ .x-title }`. The choice of delimeters ensures the YAML config file remains a valid YAML file that can be validated by the schema.

### Values

Extension fields can be reused anywhere else in
the YAML config.

```yaml
x-hello: Hello from YAML template
x-nest:
  value: nesting
```

Referencing them without `x-` prefix.

```yaml
...
handle:
  - handler: static_response
    body: "#{ .hello } with #{ .nest.value }"
```

_If string interpolation is not needed, YAML anchors and aliases can also be used to achieve this_.

### Environment Variables

Environment variables can be used in a template by prefixing with `$`.

```yaml
listen:
  - "#{ $PORT }"
...
handler: file_server
root: "#{ $APP_ROOT_DIR }/public"
```

Caddy supports runtime environment variables via [`{env.*}` placeholders](https://caddyserver.com/docs/caddyfile/concepts#environment-variables).

### Example Config 

Check the [test YAML configuration file](https://github.com/abiosoft/caddy-yaml/blob/master/testdata/test.caddy.yaml).

## License

Apache 2
