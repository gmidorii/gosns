package meta

import (
	"encoding/json"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/emluque/dscache"
	httpdoc "github.com/mercari/go-httpdoc"
	"github.com/midorigreen/gosns/topic"

	"bytes"
	"net/http"

	"path/filepath"
)

func TestSubscribeHandler(t *testing.T) {
	doc := &httpdoc.Document{
		Name: "Subscribe API",
		ExcludeHeaders: []string{
			"Accept-Encoding",
			"Content-Length",
			"User-Agent",
		},
	}
	defer func() {
		pwd, _ := os.Getwd()
		os.Setenv("HTTPDOC", "1")
		if err := doc.Generate(filepath.Join(pwd, "../doc/meta-subscribe.md")); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	pwd, _ := os.Getwd()
	path := filepath.Join(pwd, "/subscribed-test.json")
	err := createSubscribedFile(path)
	if err != nil {
		t.Error("failed test file created")
	}
	s := subscribe{
		TopicData: createTopicData(path, t),
	}

	ts, err := setTestServer("/meta/subscribe", doc, s.handler, "Register topic subscribed")
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
  "channel": "/meta/subscribe",
  "client_id": "hogehoge",
  "subscriptions" : [
  	"/golang"
  ],
  "method" : {
    "format": "slack",
    "webhook_url": "https://hooks.slack.com/services/XXXXX"
  }
}
`
	res, err := http.Post(ts.URL+"/meta/subscribe", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Error("server connection error")
	}
	defer res.Body.Close()
	var sRes subscribeRes
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&sRes)

	if sRes.Channel != "/meta/subscribe" {
		t.Error("unexpected channel name: " + sRes.Channel)
	}
	if sRes.Successful != true {
		t.Errorf("failed subscribed: %s", sRes.Error)
	}
	if sRes.ClientID != "hogehoge" {
		t.Errorf("unexpected clientID: " + sRes.ClientID)
	}

	topic, err := s.TopicData.Fetch("/golang")
	if err != nil {
		t.Errorf("failed fetch topic err: %s", err)
	}
	if topic.Channel != "/golang" {
		t.Errorf("upexpected topic from cache: %s", err)
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

func createTopicData(path string, t *testing.T) *topic.TopicData {
	ds, err := dscache.New(2 * dscache.MB)
	if err != nil {
		t.Error("failed create cache")
	}
	topics, err := topic.LoadFile(path)
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
	return &topic.TopicData{
		Path: path,
		Ds:   ds,
	}
}

func createSubscribedFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	topics := []topic.Topic{}
	topics = append(topics, topic.Topic{
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

func contains(s string, slice []topic.Subscriber) bool {
	for _, v := range slice {
		if v.ClientID == s {
			return true
		}
	}
	return false
}
