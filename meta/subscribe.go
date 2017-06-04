package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var (
	subscribedFile = "./subscribed.json"
)

type request struct {
	Channel      string   `json:"channel"`
	Subscription []string `json:"subscription"`
}

type response struct {
	Channel    string `json:"channel"`
	Successful bool   `json:"successful"`
	// ClientID     string `json:"clientId"`
	Subscription string `json:"subscription"`
	Error        string
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println(subHandler)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		res := response{
			Channel:      subHandler,
			Successful:   false,
			Subscription: "",
			Error:        "read body error",
		}
		unsuccess(res, w)
		return
	}
	defer r.Body.Close()

	req := request{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(string(body))
		log.Println(err)
		res := response{
			Channel:      subHandler,
			Successful:   false,
			Subscription: "",
			Error:        "parse request to json error",
		}
		unsuccess(res, w)
		return
	}
}

func register(req request) error {
	_, err := os.Create(subscribedFile)
	if err != nil {
		return err
	}
	return nil
}
