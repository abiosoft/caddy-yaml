package caddyyaml

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

const (
	openingDelim = "#{"
	closingDelim = "}"
)

func applyTemplate(body []byte, values map[string]interface{}, env []string) ([]byte, error) {
	tplBody := envVarsTemplate(env) + string(body)

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

func envVarsTemplate(env []string) string {
	var builder strings.Builder
	line := func(key, val string) string {
		return tplWrap(fmt.Sprintf(`$%s := %q`, key, val))
	}
	for _, env := range env {
		key, val, _ := strings.Cut(env, "=")
		fmt.Fprintln(&builder, line(key, val))
	}
	return builder.String()
}

func tplWrap(s string) string {
	return fmt.Sprintf("%s %s %s", openingDelim, s, closingDelim)
}
