package caddyyaml

import (
	"io/ioutil"
	"testing"
)

type M map[string]interface{}

func TestApply(t *testing.T) {
	b, err := ioutil.ReadFile("./testdata/test.caddy.yaml")
	if err != nil {
		t.Fatal(err)
	}

	b, _, err = Adapter{}.Adapt(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}
