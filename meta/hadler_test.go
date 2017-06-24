package meta

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

type Sample struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestDecodeBody(t *testing.T) {
	json := `{
		"id" : 1,
		"name" : "hoge"
}`
	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(json))
	var s Sample
	decodeBody(req, &s)

	if s.ID != 1 {
		t.Error("failed id")
		t.Log("s.ID: " + string(s.ID) + ", id: 1")
	}
	if s.Name != "hoge" {
		t.Error("failed name")
		t.Log("s.Name: " + s.Name + ", Name: hoge")
	}
}
