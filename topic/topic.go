package topic

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
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
	FilePath = "./subscribed.json"
)

var BufTopics = sync.Pool{
	New: func() interface{} {
		topics := []Topic{}
		return &topics
	},
}

func init() {
	BufTopics.Put(LoadTopics())
}

// LoadTopics is loading file registered topic information
func LoadTopics() *[]Topic {
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
	return &s
}

// ReLoadTopics is topics cache reloading
func ReLoadTopics() {
	BufTopics.Get()
	topics := LoadTopics()
	BufTopics.Put(topics)
}
