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
	ClientID     string   `json:"client_id"`
	Subscription []string `json:"subscription"`
}

type response struct {
	Channel      string   `json:"channel"`
	Successful   bool     `json:"successful"`
	ClientID     string   `json:"clientId"`
	Subscription []string `json:"subscription"`
	Error        string   `json:",omitempty"`
}

type subscribed struct {
	Channel string
	Clients []string
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

	if err = register(req); err != nil {
		log.Println(err)
		unsuccessed("register faild", "", []string{}, w)
		return
	}

	successed(req.ClientID, req.Subscription, w)
}

func register(req request) error {
	file, err := os.OpenFile(subscribedFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	v, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	s := []subscribed{}
	err = json.Unmarshal(v, &s)
	if err != nil {
		return err
	}
	if len(s) != 0 {
		for _, v := range s {
			if contains(req.Subscription, v.Channel) {
				if !contains(v.Clients, req.ClientID) {
					v.Clients = append(v.Clients, req.ClientID)
				}
			} else {
				newSub := subscribed{
					Channel: req.Channel,
					Clients: []string{req.ClientID},
				}
				s = append(s, newSub)
			}
		}
	} else {
		s = append(s, subscribed{
			Channel: req.Channel,
			Clients: []string{req.ClientID},
		})
	}

	byte, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(subscribedFile, byte, 0666)
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
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
