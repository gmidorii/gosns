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

type subscribed struct {
	Channel    string
	Subscribes []string
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
		writeRes(res, w)
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
		writeRes(res, w)
		return
	}

	if err = register(req); err != nil {
		log.Println(err)
		res := response{
			Channel:      subHandler,
			Successful:   false,
			Subscription: "",
			Error:        "regiser faild",
		}
		writeRes(res, w)
		return
	}
	res := response{
		Channel:      subHandler,
		Successful:   true,
		Subscription: "",
	}
	writeRes(res, w)
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
			if v.Channel == req.Channel {
				for _, sub := range req.Subscription {
					if !contains(v.Subscribes, sub) {
						v.Subscribes = append(v.Subscribes, sub)
					}
				}
			} else {
				v.Subscribes = append(v.Subscribes, req.Subscription...)
			}
		}
	} else {
		s = append(s, subscribed{
			Channel:    req.Channel,
			Subscribes: req.Subscription,
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
