package channel

import (
	"encoding/json"
	"net/http"
)

const (
	topicsPath = "/topics"
)

type topicReq struct {
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// Handler is topic handler
func Handler() {
	http.HandleFunc(topicsPath, topicsHandler)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	var tReq topicReq
	decodeBody(r, &tReq)
	topic := findChannel(PoolTopics.Get().([]Topic), tReq.Channel)
	if topic.Channel == "" {
		w.Write([]byte("not found channel"))
		return
	}
	go sends(topic, tReq)
	w.Write([]byte("ok"))
}

func decodeBody(req *http.Request, out interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(out)
}

func findChannel(topics []Topic, channel string) Topic {
	for _, v := range topics {
		if channel == v.Channel {
			return v
		}
	}
	return Topic{}
}
