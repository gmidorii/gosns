package meta

import (
	"encoding/json"
	"net/http"

	"github.com/midorigreen/gosns/topic"
)

var (
	handshakePattarn = "/meta/handshake"
	subscribePattarn = "/meta/subscribe"
	topicPattarn     = "/meta/topic"
)

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(handshakePattarn, handshakeHandler)
	s := subscribe{
		TopicData: topic.CreateTopicData(),
	}
	http.HandleFunc(subscribePattarn, s.handler)
	c := channel{
		TopicData: topic.CreateTopicData(),
	}
	http.HandleFunc(topicPattarn, c.handler)
}

func decodeBody(req *http.Request, out interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(out)
}

func writeRes(v interface{}, w http.ResponseWriter, statusCode int) {
	json, err := json.Marshal(v)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(json)
}
