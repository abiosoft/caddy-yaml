package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"strings"
	"unicode"
)

var variableLineRegexp = regexp.MustCompile(`^\_([a-zA-Z0-9]+)(\s*)\:`)

func commentLine(line string) bool { return strings.HasPrefix(strings.TrimSpace(line), "#") }
func topLevelLine(line string) bool {
	return strings.IndexFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r)
	}) == 0 && !commentLine(line) && line[0] != '_'
}
func variableLine(line string) (variable string, found bool) {
	// TODO: use other than regexp if this is notably slow.
	if matched := variableLineRegexp.FindStringSubmatch(line); len(matched) > 0 {
		found = true
		variable = matched[1]
	}
	return
}

// extractVariables extracts variables from the body.
func extractVariables(body []byte) ([]byte, error) {
	var variablesBuffer bytes.Buffer

	reader := bufio.NewReader(bytes.NewReader(body))
	for {
		err := bufReadUntil(reader, func(line string) bool {
			_, ok := variableLine(line)
			if !ok {
				return false
			}
			variablesBuffer.WriteString(line)
			return true
		})
		// discard EOF as error
		if err != nil && err != io.EOF {
			return nil, err
		}

		err = bufReadUntil(reader, func(line string) bool {
			if !topLevelLine(line) {
				variablesBuffer.WriteString(line)
				return false
			}
			return true
		})
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return variablesBuffer.Bytes(), nil
}

func bufReadUntil(buf *bufio.Reader, f func(string) bool) error {
	var line string
	var err error

	for {
		line, err = buf.ReadString('\n')
		if err != nil || f(line) {
			break
		}
	}

	return err
}
