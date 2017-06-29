package meta

import (
	"bytes"
	"testing"

	"path/filepath"

	"net/http"

	"net/http/httptest"

	httpdoc "github.com/mercari/go-httpdoc"
)

func TestHandler(t *testing.T) {
	doc := &httpdoc.Document{
		Name:           "Handshake API",
		ExcludeHeaders: []string{},
	}
	defer func() {
		if err := doc.Generate(filepath.Join("doc", "handshake.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/meta/handshake", httpdoc.Record(http.HandlerFunc(handshakeHandler), doc, &httpdoc.RecordOption{Description: "Handshake Server"}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	reqBody := `
{
	"channel": "/meta/handshake"
}
`

	res, err := http.Post(ts.URL+"/meta/handshake", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Error("server connection error")
	}
}
