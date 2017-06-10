package topic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

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

const (
	// FilePath is file path registerd topics information
	FilePath   = "./subscribed.json"
	topicsPath = "/topics"
)

var topics = LoadTopics()

// Handler is topic handler
func Handler() {
	http.HandleFunc(topicsPath, topicsHandler)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(topics)
}

// LoadTopics is loading topics in file
// You must execute this when add topic or subscriber
func LoadTopics() []Topic {
	file, err := os.OpenFile(FilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	defer file.Close()

	v, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	s := []Topic{}
	err = json.Unmarshal(v, &s)
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return s
}
