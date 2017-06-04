package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	subHandler = "/meta/subscribe"
)

type subscriptionReq struct {
	Channels []string
}

type subscriptionRes struct {
	Channel    string `json:"channel"`
	Successful bool   `json:"successful"`
	// ClientID     string `json:"clientId"`
	Subscription string `json:"subscription"`
}

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(subHandler, subscribe)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println(subHandler)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		unsuccess(subHandler, "", w)
		return
	}

	subReq, err := json.Marshal(body)
	if err != nil {
		unsuccess(subHandler, "", w)
		return
	}
}

func unsuccess(channel string, subscritpion string, w http.ResponseWriter) {
	res := subscriptionRes{
		Channel:      subHandler,
		Successful:   false,
		Subscription: "",
	}
	json, err := json.Marshal(res)
	if err != nil {
		return
	}
	w.Write(json)
}
