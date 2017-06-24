package meta

import (
	"errors"
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/channel"
)

type subscribeReq struct {
	Channel       string    `json:"channel"`
	ClientID      string    `json:"client_id"`
	Subscriptions []string  `json:"subscription"`
	Method        strMethod `json:"method"`
}

type strMethod struct {
	Format     string `json:"format"`
	WebHookURL string `json:"webhook_url"`
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

func register(sReq subscribeReq) ([]string, error) {
	topics := *channel.LoadTopics()
	if len(topics) == 0 {
		return nil, errors.New("not found topic")
	}

	var registered = []string{}
	for i, v := range topics {
		if checkAppendSubscribed(sReq, v) {
			format := channel.FormatValue(sReq.Method.Format)
			if format == channel.Error {
				continue
			}
			topics[i].Subscribers = append(topics[i].Subscribers, channel.Subscriber{
				ClientID: sReq.ClientID,
				Method: channel.Method{
					Format:     format,
					WebFookURL: sReq.Method.WebHookURL,
				},
			})
			registered = append(registered, v.Channel)
		}
	}
	if len(registered) == 0 {
		return nil, errors.New("not founc topic")
	}

	if err := channel.PutTopics(topics); err != nil {
		return nil, errors.New("register subscribed failed")
	}
	return registered, nil
}

func checkAppendSubscribed(sReq subscribeReq, t channel.Topic) bool {
	return containsSlice(sReq.Subscriptions, t.Channel) && !containsSubscriber(sReq.ClientID, t.Subscribers)
}

func containsSlice(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func containsSubscriber(e string, s []channel.Subscriber) bool {
	for _, v := range s {
		if e == v.ClientID {
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
	res := subscribeRes{
		Channel:      subHandler,
		Successful:   success,
		ClientID:     clientID,
		Subscription: sub,
		Error:        mes,
	}
	writeRes(res, w)
}
