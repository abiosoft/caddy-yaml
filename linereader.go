package caddyyaml

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

func commentLine(line string) bool { return strings.HasPrefix(strings.TrimSpace(line), "#") }
func topLevelLine(line string) bool {
	return strings.IndexFunc(line, func(r rune) bool {
		return !unicode.IsSpace(r) && !commentLine(line)
	}) == 0
}
func configLine(line string, configs ...string) (config string, found bool) {
	for _, config = range configs {
		if strings.HasPrefix(line, config) {
			line = strings.TrimPrefix(line, config)
			if strings.HasPrefix(strings.TrimSpace(line), ":") {
				found = true
				break
			}
		}
	}
	return
}

// extractConfigs extracts configs from the body. Returns the body, config, error.
func extractConfigs(body []byte, configs ...string) ([]byte, []byte, error) {
	set := map[string]struct{}{}
	for _, config := range configs {
		set[config] = struct{}{}
	}

	var configBuffer bytes.Buffer
	var bodyBuffer bytes.Buffer

	reader := bufio.NewReader(bytes.NewReader(body))
	for len(set) > 0 {
		err := bufReadUntil(reader, func(line string) bool {
			config, ok := configLine(line, configs...)
			if !ok {
				bodyBuffer.WriteString(line)
				return false
			}
			delete(set, config)
			configBuffer.WriteString(line)
			return true
		})
		// discard EOF as error
		if err != nil && err != io.EOF {
			return nil, nil, err
		}

		err = bufReadUntil(reader, func(line string) bool {
			if !topLevelLine(line) {
				configBuffer.WriteString(line)
				return false
			}
			bodyBuffer.WriteString(line)
			return true
		})
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, nil, err
		}
	}

	io.Copy(&bodyBuffer, reader)

	return bodyBuffer.Bytes(), configBuffer.Bytes(), nil
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
