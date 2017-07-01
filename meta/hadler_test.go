package meta

import (
	"bytes"
	"fmt"
	"net/http"
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
		t.Log("s.ID: " + fmt.Sprint(s.ID) + ", id: 1")
	}
	if s.Name != "hoge" {
		t.Error("failed name")
		t.Log("s.Name: " + s.Name + ", Name: hoge")
	}
}

func TestWriteRes(t *testing.T) {
	w := httptest.NewRecorder()
	sample := Sample{2, "hoge"}
	writeRes(sample, w, http.StatusOK)
	resp := w.Result()
	if len(resp.Header) != 1 {
		t.Error("unexpected header length")
		t.Log(len(resp.Header))
	}
	for key, v := range resp.Header {
		if key != "Content-Type" {
			t.Error("unexpected header key")
			t.Log(key)
		}
		if len(v) != 1 {
			t.Error("unexpected value length")
			t.Log(len(v))
		}
		if v[0] != "application/json" {
			t.Error("unexpected header value")
			t.Log(v)
		}
	}
}
