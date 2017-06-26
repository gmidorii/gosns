package meta

import (
	"net/http"

	"github.com/midorigreen/gopubsub/channel"
)

type topic struct {
	TopicData *channel.TopicData
}

type topicReq struct {
	Channel string `json:"channel"`
}

type topicRes struct {
	Channel    string `json:"channel,omitempty"`
	Successful bool   `json:"successful"`
	Error      string `json:"error,omitempty"`
}

func (t *topic) handler(w http.ResponseWriter, r *http.Request) {
	var tReq topicReq
	if err := decodeBody(r, &tReq); err != nil {
		writeRes(topicRes{
			Successful: false,
			Error:      "failed parsing request body",
		}, w)
	}

	topic := channel.Topic{
		Channel: tReq.Channel,
	}

	if err := t.TopicData.Add(topic); err != nil {
		writeRes(topicRes{
			Channel:    tReq.Channel,
			Successful: false,
			Error:      err.Error(),
		}, w)
	}
	writeRes(topicRes{
		Channel:    tReq.Channel,
		Successful: true,
	}, w)
}
