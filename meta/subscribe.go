package meta

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/topic"
)

type subscribeReq struct {
	Channel       string   `json:"channel"`
	ClientID      string   `json:"client_id"`
	Subscriptions []string `json:"subscription"`
}

type subscribeRes struct {
	Channel      string   `json:"channel"`
	Successful   bool     `json:"successful"`
	ClientID     string   `json:"clientId"`
	Subscription []string `json:"subscription"`
	Error        string   `json:",omitempty"`
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println(subHandler)
	req := subscribeReq{}
	err := decodeBody(r, &req)
	if err != nil {
		log.Println(err)
		unsuccessed("parse request to json error", "", []string{}, w)
		return
	}

	registered, err := register(req)
	if err != nil {
		log.Println(err)
		unsuccessed(err.Error(), req.ClientID, []string{}, w)
		return
	}

	successed(req.ClientID, registered, w)
}

func register(req subscribeReq) ([]string, error) {
	s := *topic.LoadTopics()
	if len(s) == 0 {
		return nil, errors.New("not found topic")
	}

	var registered = []string{}
	for i, v := range s {
		if containsSlice(v.Channel, req.Subscriptions) {
			if !containsSubscriber(req.ClientID, v.Subscribers) {
				s[i].Subscribers = append(s[i].Subscribers, topic.Subscriber{
					ClientID: req.ClientID,
				})

			}
			registered = append(registered, v.Channel)
		}
	}
	if len(registered) == 0 {
		return nil, errors.New("not founc topic")
	}

	if err := writeSubscribed(s); err != nil {
		return nil, errors.New("register subscribed failed")
	}
	return registered, nil
}

func containsSlice(e string, s []string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func containsSubscriber(e string, s []topic.Subscriber) bool {
	for _, v := range s {
		if e == v.ClientID {
			return true
		}
	}
	return false
}

func writeSubscribed(s []topic.Topic) error {
	byte, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(topic.FilePath, byte, 0666); err != nil {
		return err
	}
	topic.ReLoadTopics()
	return nil
}

func unsuccessed(errMes, clientID string, sub []string, w http.ResponseWriter) {
	write(false, errMes, clientID, sub, w)
}

func successed(clientID string, sub []string, w http.ResponseWriter) {
	write(true, "", clientID, sub, w)
}

func write(success bool, mes, clientID string, sub []string, w http.ResponseWriter) {
	res := subscribeRes{
		Channel:      subHandler,
		Successful:   success,
		ClientID:     clientID,
		Subscription: sub,
		Error:        mes,
	}
	writeRes(res, w)
}
