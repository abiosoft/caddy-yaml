package caddyyaml

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestApply(t *testing.T) {
	tests := []struct {
		filename    string
		environment string
	}{
		{filename: "test.caddy.json"},
		{filename: "test.caddy.prod.json", environment: "production"},
	}

	for _, tt := range tests {
		t.Run(tt.environment, func(t *testing.T) {
			os.Setenv("ENVIRONMENT", tt.environment)

			b, err := ioutil.ReadFile("./testdata/test.caddy.yaml")
			if err != nil {
				t.Fatal(err)
			}
			adaptedBytes, _, err := Adapter{}.Adapt(b, nil)
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
