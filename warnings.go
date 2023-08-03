package caddyyaml

import "github.com/caddyserver/caddy/v2/caddyconfig"

type warningsCollector struct {
	filename string

	warnings []caddyconfig.Warning
}

func (w *warningsCollector) Add(line int, directive string, message string) {
	w.warnings = append(w.warnings, caddyconfig.Warning{
		File:      w.filename,
		Line:      line,
		Directive: directive,
		Message:   message,
	})
}

func newWarningsCollector(filename string) *warningsCollector {
	return &warningsCollector{filename, nil}
}
