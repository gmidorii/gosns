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
	var t topicReq
	decodeBody(r, &t)
	topic := findChannel(GetTopics(), t.Channel)
	if topic.Channel == "" {
		w.Write([]byte("not found channel"))
		return
	}
	go func(topic Topic, t topicReq) {
		for _, v := range topic.Subscribers {
			var service SubscriberService
			switch v.Method.Format {
			case Slack:
				service = &SlackSender{
					URL: v.Method.WebFookURL,
				}
			// case Mail:
			default:
				continue
			}
			service.Send(t.Data)
		}
	}(topic, t)
	w.Write([]byte(t.Channel))
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
