package meta

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"bytes"

	httpdoc "github.com/mercari/go-httpdoc"
)

func TestTopicHandler(t *testing.T) {
	doc := &httpdoc.Document{
		Name: "Topic append API",
		ExcludeHeaders: []string{
			"Accept-Encoding",
			"Content-Length",
			"User-Agent",
		},
	}
	defer func() {
		pwd, _ := os.Getwd()
		os.Setenv("HTTPDOC", "1")
		if err := doc.Generate(filepath.Join(pwd, "../doc/meta-topic.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "subscribed-test.json")
	err := createSubscribedFile(path)
	if err != nil {
		t.Error("failed test file created")
	}
	topic := topic{
		TopicData: createTopicData(path, t),
	}
	mux := http.NewServeMux()
	apiPath := "/meta/topic"
	mux.Handle(apiPath, httpdoc.Record(http.HandlerFunc(topic.handler), doc, &httpdoc.RecordOption{Description: "Register new topic"}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	reqBody := `
{
	"channel": "govim"
}
`

	res, err := http.Post(ts.URL+apiPath, "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Error("server connection error")
	}
	defer res.Body.Close()
	var topicRes topicRes
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&topicRes)

	if topicRes.Channel != "govim" {
		t.Errorf("unexpected channel name: %s", topicRes.Channel)
	}
	if topicRes.Successful != true {
		t.Errorf("failed topic append: %s", topicRes.Error)
	}

	// tear down
	if err = deleteTestFile(path); err != nil {
		t.Error("failed delete file")
	}
}
