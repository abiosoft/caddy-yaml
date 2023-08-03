package caddyyaml

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name             string
		filename         string
		env              []string
		expectedWarnings []string
	}{
		{
			name:     "simple",
			filename: "test.caddy.json",
			env:      []string{"ENVIRONMENT=something"},
		},
		{
			name:     "production environment",
			filename: "test.caddy.prod.json",
			env:      []string{"ENVIRONMENT=production"},
		},
		{
			name:     "bad environment variables",
			filename: "test.caddy.json",
			env: []string{
				"ENVIRONMENT=bad",
				"INVALID%=invalid_name",
				"VALUE_WITH_BACKSLASH=foo\\bar",
			},
			expectedWarnings: []string{
				"test.caddy.yaml:-1: environment variable \"INVALID%\" cannot be used in template",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ioutil.ReadFile("./testdata/test.caddy.yaml")
			if err != nil {
				t.Fatal(err)
			}
			adaptedBytes, warnings, err := Adapter{}.Adapt(b, map[string]interface{}{
				"filename":    "test.caddy.yaml",
				envOptionName: tt.env,
			})
			if err != nil {
				t.Fatal(err)
			}

			for i, w := range warnings {
				if len(tt.expectedWarnings) < i+1 {
					t.Fatalf("unexpected warning %q", w)
				}
				eW := tt.expectedWarnings[i]
				if eW != w.String() {
					t.Fatalf("expected warning %q, got %q", eW, w)
				}
			}
			if len(tt.expectedWarnings) > len(warnings) {
				t.Fatalf("expected additional warnings: %v", tt.expectedWarnings[len(warnings):])
			}

			jsonBytes, err := ioutil.ReadFile("./testdata/" + tt.filename)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(jsonToObj(adaptedBytes), jsonToObj(jsonBytes)) {
				t.Log(string(adaptedBytes))
				t.Fatal("adapter config does not match expected config")
			}

		})
	}

}

func jsonToObj(b []byte) (obj map[string]interface{}) {
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
	return
}
