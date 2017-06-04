package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
