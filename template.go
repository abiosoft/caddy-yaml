package caddyyaml

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const (
	openingDelim = "#{"
	closingDelim = "}"
)

func applyTemplate(body []byte, values map[string]interface{}, env []string, wc *warningsCollector) ([]byte, error) {
	tplBody := envVarsTemplate(env, wc) + string(body)

	tpl, err := template.New("yaml").
		Funcs(sprig.TxtFuncMap()).
		Delims(openingDelim, closingDelim).
		Parse(tplBody)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, values); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func envVarsTemplate(env []string, wc *warningsCollector) string {
	var builder strings.Builder
	line := func(key, val string) string {
		return tplWrap(fmt.Sprintf(`$%s := %q`, key, val))
	}
	for _, env := range env {
		key, val, _ := strings.Cut(env, "=")
		if !token.IsIdentifier(key) {
			if wc != nil {
				wc.Add(-1, "", fmt.Sprintf("environment variable %q cannot be used in template", key))
			}
			continue
		}
		fmt.Fprintln(&builder, line(key, val))
	}
	return builder.String()
}

func tplWrap(s string) string {
	return fmt.Sprintf("%s %s %s", openingDelim, s, closingDelim)
}
