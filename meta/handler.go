package meta

import (
	"encoding/json"
	"net/http"
)

var (
	subHandler = "/meta/subscribe"
)

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(subHandler, subscribe)
}

func unsuccess(channel string, subscritpion string, w http.ResponseWriter) {
	res := subscriptionRes{
		Channel:      channel,
		Successful:   false,
		Subscription: subscritpion,
	}
	json, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(json)
}
