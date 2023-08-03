package caddyyaml

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		env      []string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ioutil.ReadFile("./testdata/test.caddy.yaml")
			if err != nil {
				t.Fatal(err)
			}
			adaptedBytes, _, err := Adapter{}.Adapt(b, map[string]interface{}{
				"filename":    "test.caddy.yaml",
				envOptionName: tt.env,
			})
			if err != nil {
				t.Fatal(err)
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
