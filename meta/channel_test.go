package meta

import (
	"encoding/json"
	"net/http"
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
		if err := doc.Generate(filepath.Join(pwd, "../doc/meta-channel.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	apiPath := "/meta/channel"
	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "subscribed-test.json")
	err := createSubscribedFile(path)
	if err != nil {
		t.Errorf("failed create subscribed file: %s", err)
	}
	topic := channel{
		TopicData: createTopicData(path, t),
	}
	ts, err := setTestServer(apiPath, doc, topic.handler, "Register new topic")
	defer func() {
		ts.Close()
		if err := deleteTestFile(path); err != nil {
			t.Error("failed delete file")
		}
	}()

	if err != nil {
		t.Errorf("failed setting test server: %s", err)
	}

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
	var channelRes channelRes
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&channelRes)

	if channelRes.Channel != "govim" {
		t.Errorf("unexpected channel name: %s", channelRes.Channel)
	}
	if channelRes.Successful != true {
		t.Errorf("failed topic append: %s", channelRes.Error)
	}

}
