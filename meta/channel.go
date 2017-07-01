package meta

import (
	"net/http"

	"github.com/midorigreen/gosns/topic"
)

type channel struct {
	TopicData *topic.TopicData
}

type channelReq struct {
	Channel string `json:"channel"`
}

type channelRes struct {
	Channel    string `json:"channel,omitempty"`
	Successful bool   `json:"successful"`
	Error      string `json:"error,omitempty"`
}

func (t *channel) handler(w http.ResponseWriter, r *http.Request) {
	var tReq channelReq
	if err := decodeBody(r, &tReq); err != nil {
		writeRes(channelRes{
			Successful: false,
			Error:      "failed parsing request body",
		}, w, http.StatusBadRequest)
	}

	topic := topic.Topic{
		Channel: tReq.Channel,
	}

	if err := t.TopicData.Add(topic); err != nil {
		writeRes(channelRes{
			Channel:    tReq.Channel,
			Successful: false,
			Error:      err.Error(),
		}, w, http.StatusInternalServerError)
	}
	writeRes(channelRes{
		Channel:    tReq.Channel,
		Successful: true,
	}, w, http.StatusOK)
}
