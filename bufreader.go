package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

func commentLine(line string) bool { return strings.HasPrefix(strings.TrimSpace(line), "#") }
func emptyLine(line string) bool   { return strings.TrimSpace(line) == "#" }
func topLevelLine(line string) bool {
	return strings.IndexFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r) && !commentLine(line)
	}) == 0
}
func hasMapKey(line string, keys ...string) (key string, found bool) {
	for _, key = range keys {
		if strings.HasPrefix(line, key) {
			line = strings.TrimPrefix(line, key)
			if strings.HasPrefix(strings.TrimSpace(line), ":") {
				found = true
				break
			}
		}
	}
	return
}

// extractConfigs extracts configs from the body. Returns the body, config, error.
func extractConfigs(body []byte, keys ...string) ([]byte, []byte, error) {
	keySet := map[string]struct{}{}
	for _, key := range keys {
		keySet[key] = struct{}{}
	}

	var cout bytes.Buffer
	var bout bytes.Buffer

	reader := bufio.NewReader(bytes.NewReader(body))
	for len(keySet) > 0 {
		err := bufReadUntil(reader, func(line string) bool {
			key, ok := hasMapKey(line, keys...)
			if !ok {
				bout.WriteString(line)
				return false
			}
			delete(keySet, key)
			cout.WriteString(line)
			return true
		})
		// discard EOF as error
		if err != nil && err != io.EOF {
			return nil, nil, err
		}

		err = bufReadUntil(reader, func(line string) bool {
			if !topLevelLine(line) {
				cout.WriteString(line)
				return false
			}
			bout.WriteString(line)
			return true
		})
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
	}

	io.Copy(&bout, reader)

	return bout.Bytes(), cout.Bytes(), nil
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
