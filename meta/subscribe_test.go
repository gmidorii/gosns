package meta

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	httpdoc "github.com/mercari/go-httpdoc"
	"github.com/midorigreen/gopubsub/channel"

	"bytes"
	"net/http"
)

func createTopicPool(path string, t *testing.T) *channel.TopicPool {
	return &channel.TopicPool{
		Value: sync.Pool{
			New: func() interface{} {
				topics, err := channel.LoadFile(path)
				if err != nil {
					t.Error("failed subscribed file read")
				}
				return &topics
			},
		},
		Path: path,
	}
}

func createSubscribedFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	topics := []channel.Topic{}
	topics = append(topics, channel.Topic{
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

func contains(s string, slice []channel.Subscriber) bool {
	for _, v := range slice {
		if v.ClientID == s {
			return true
		}
	}
	return false
}

func TestHandler(t *testing.T) {
	document := &httpdoc.Document{
		Name:           "Subscribe API",
		ExcludeHeaders: []string{},
	}
	defer func() {
		if err := document.Generate("doc/subscribe.md"); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	path := "./subscribed-test.json"
	err := createSubscribedFile(path)
	if err != nil {
		t.Error("failed test file created")
	}
	s := Subscribe{
		TopicPool: createTopicPool(path, t),
	}
	mux := http.NewServeMux()
	mux.Handle("/meta/subscribe", httpdoc.Record(http.HandlerFunc(s.handler), document, &httpdoc.RecordOption{Description: "Register topic subscribed"}))
	ts := httptest.NewServer(mux)
	defer ts.Close()

	reqBody := `
{
  "channel": "/meta/subscribe",
  "client_id": "hogehoge",
  "subscription" : [
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
		t.Error("failed subscribed")
	}
	if sRes.ClientID != "hogehoge" {
		t.Error("unexpected clientID: " + sRes.ClientID)
	}

	// cache delete (temporary solution)
	s.TopicPool.Value.Get()
	topics := s.TopicPool.Get().([]channel.Topic)
	f := false
	for _, v := range topics {
		if v.Channel == "/golang" {
			f = contains("hogehoge", v.Subscribers)
			break
		}
	}
	if f != true {
		t.Error("failed cache clear pool")
	}

	if err = deleteTestFile(path); err != nil {
		t.Error("failed delete file")
	}
}
