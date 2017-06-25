package meta

import (
	"encoding/json"
	"net/http"

	"github.com/midorigreen/gopubsub/channel"
)

var (
	handshakePattarn = "/meta/handshake"
	subscribePattarn = "/meta/subscribe"
	topicPattarn     = "/meta/topic"
)

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(handshakePattarn, handshakeHandler)
	s := Subscribe{
		TopicPool: &channel.PoolTopics,
	}
	http.HandleFunc(subscribePattarn, s.handler)
	http.HandleFunc(topicPattarn, topicHandler)
}

func decodeBody(req *http.Request, out interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(out)
}

func writeRes(v interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(v)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
