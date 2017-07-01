package meta

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"path/filepath"

	"net/http"

	httpdoc "github.com/mercari/go-httpdoc"
)

func TestHandshakeHandler(t *testing.T) {
	doc := &httpdoc.Document{
		Name: "Handshake API",
		ExcludeHeaders: []string{
			"Accept-Encoding",
			"Content-Length",
			"User-Agent",
		},
	}
	defer func() {
		pwd, _ := os.Getwd()
		os.Setenv("HTTPDOC", "1")
		if err := doc.Generate(filepath.Join(pwd, "../doc/meta-handshake.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	ts, err := setTestServer("/meta/handshake", doc, handshakeHandler, "Handshake server")
	if err != nil {
		t.Errorf("failed create test server: %s", err)
	}
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
	defer res.Body.Close()
	var hRes handShakeRes
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&hRes)

	if hRes.Channel != "/meta/handshake" {
		t.Errorf("unexpected channel name: %s", hRes.Channel)
	}
	if hRes.Successful != true {
		t.Errorf("failed handshake")
	}
	if len(hRes.ClientID) != 10 {
		t.Errorf("unexpected clientid length: %d", len(hRes.ClientID))
	}
}
