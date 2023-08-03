package caddyyaml

import (
	"bytes"
	"fmt"
	"strconv"
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
		return body, err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, values); err != nil {
		return body, err
	}
	return out.Bytes(), nil
}

func envVarsTemplate(env []string) string {
	var builder strings.Builder
	line := func(key, val string) string {
		// avoid quoted string
		if len(val) > 0 && val[0] == '"' {
			if v, err := strconv.Unquote(val); err == nil {
				val = v
			}
		}
		return tplWrap(fmt.Sprintf(`$%s := "%s"`, key, val))
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
