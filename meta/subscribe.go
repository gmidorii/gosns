package meta

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/midorigreen/gopubsub/topic"
)

var (
	subscribedFile = "./subscribed.json"
)

type request struct {
	Channel       string   `json:"channel"`
	ClientID      string   `json:"client_id"`
	Subscriptions []string `json:"subscription"`
}

type response struct {
	Channel      string   `json:"channel"`
	Successful   bool     `json:"successful"`
	ClientID     string   `json:"clientId"`
	Subscription []string `json:"subscription"`
	Error        string   `json:",omitempty"`
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println(subHandler)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		unsuccessed("read body error", "", []string{}, w)
		return
	}
	defer r.Body.Close()

	req := request{}
	err = json.Unmarshal(body, &req)
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

func register(req request) ([]string, error) {
	s, err := readSubscribed()
	if err != nil {
		return nil, err
	}
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

	if err = writeSubscribed(s); err != nil {
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

func readSubscribed() ([]topic.Topic, error) {
	file, err := os.OpenFile(subscribedFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	v, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	s := []topic.Topic{}
	err = json.Unmarshal(v, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func writeSubscribed(s []topic.Topic) error {
	byte, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(subscribedFile, byte, 0666)
}

func unsuccessed(errMes, clientID string, sub []string, w http.ResponseWriter) {
	write(false, errMes, clientID, sub, w)
}

func successed(clientID string, sub []string, w http.ResponseWriter) {
	write(true, "", clientID, sub, w)
}

func write(success bool, mes, clientID string, sub []string, w http.ResponseWriter) {
	res := response{
		Channel:      subHandler,
		Successful:   success,
		ClientID:     clientID,
		Subscription: sub,
		Error:        mes,
	}
	writeRes(res, w)
}
