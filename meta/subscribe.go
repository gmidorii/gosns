package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type request struct {
	Channels []string
}

type response struct {
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

	req := request{}
	err = json.Unmarshal(body, req)
	if err != nil {
		unsuccess(subHandler, "", w)
		return
	}
}
