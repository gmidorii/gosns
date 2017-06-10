package topic

import "net/http"

// Topic is main struct in package topic
// This has channle name and subscriber list
type Topic struct {
	Channel     string
	Subscribers []Subscriber
}

// Subscriber has subscriber imformation only
type Subscriber struct {
	ClientID string
}

const topicsPath = "/topics"

var topics []Topic

func init() {

}

// Handler is topic handler
func Handler() {
	http.HandleFunc(topicsPath, topicsHandler)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
}

func readTopic() {

}
