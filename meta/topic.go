package meta

import (
	"net/http"

	"github.com/midorigreen/gopubsub/channel"
)

type topicReq struct {
	Channel string `json:"channel"`
}

type topicRes struct {
	Channel    string `json:"channel,omitempty"`
	Successful bool   `json:"successful"`
	Error      string `json:"error,omitempty"`
}

func topicHandler(w http.ResponseWriter, r *http.Request) {
	topics := channel.PoolTopics.Get().([]channel.Topic)
	if len(topics) == 0 {
		topics = []channel.Topic{}
	}

	var tReq topicReq
	if err := decodeBody(r, &tReq); err != nil {
		writeRes(topicRes{
			Successful: false,
			Error:      "failed parsing request body",
		}, w)
	}

	topics = append(topics, channel.Topic{
		Channel: tReq.Channel,
	})

	if err := channel.PoolTopics.Put(topics); err != nil {
		writeRes(topicRes{
			Channel:    tReq.Channel,
			Successful: false,
			Error:      "failed append new topic",
		}, w)
	}
	writeRes(topicRes{
		Channel:    tReq.Channel,
		Successful: true,
	}, w)
}
