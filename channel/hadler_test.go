package channel

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/emluque/dscache"
	httpdoc "github.com/mercari/go-httpdoc"
)

func TestHandler(t *testing.T) {
	doc := &httpdoc.Document{
		Name: "Publishe topic API",
		ExcludeHeaders: []string{
			"Accept-Encoding",
			"Content-Length",
			"User-Agent",
		},
	}
	defer func() {
		pwd, _ := os.Getwd()
		os.Setenv("HTTPDOC", "1")
		if err := doc.Generate(filepath.Join(pwd, "../doc/topic.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "/subscribed-test.json")
	err := createSubscribedFile(path)
	if err != nil {
		t.Error("failed test file created")
	}
	topic := topic{
		TopicData: createTopicData(path, t),
	}

	ts, err := setTestServer("/topic", doc, topic.handler, "Publish topic")
	if err != nil {
		t.Errorf("failed create test server: %s", err)
	}
	defer func() {
		ts.Close()
		// tear down
		if err = deleteTestFile(path); err != nil {
			t.Error("failed delete file")
		}
	}()

	reqBody := `
{
	"channel": "/golang",
	"data": "GAE go1.8"
}
`
	res, err := http.Post(ts.URL+"/topic", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Error("server connection error")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("failed read body: %s", err)
	}
	if string(body) != "ok" {
		t.Errorf("unexpected response body: %s", string(body))
	}
}

/**
 * support function
 */

func setTestServer(apiPath string, doc *httpdoc.Document, handler func(http.ResponseWriter, *http.Request), docDescription string) (*httptest.Server, error) {
	mux := http.NewServeMux()
	mux.Handle(apiPath, httpdoc.Record(http.HandlerFunc(handler), doc, &httpdoc.RecordOption{Description: docDescription}))
	ts := httptest.NewServer(mux)
	return ts, nil
}

func createTopicData(path string, t *testing.T) *TopicData {
	ds, err := dscache.New(2 * dscache.MB)
	if err != nil {
		t.Error("failed create cache")
	}
	topics, err := LoadFile(path)
	if err != nil {
		t.Error("failed load topics")
	}

	// set cache
	for _, v := range topics {
		subByte, err := json.Marshal(v.Subscribers)
		if err != nil {
			log.Fatalln(err)
		}
		if err = ds.Set(v.Channel, string(subByte), 24*time.Hour); err != nil {
			log.Fatalln(err)
		}
	}
	return &TopicData{
		Path: path,
		Ds:   ds,
	}
}

func createSubscribedFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	topics := []Topic{}
	topics = append(topics, Topic{
		Channel: "/golang",
	})
	jsonTopics, err := json.Marshal(topics)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(jsonTopics))
	return err
}

func deleteTestFile(path string) error {
	return os.Remove(path)
}
