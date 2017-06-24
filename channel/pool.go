package channel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	poolPath = "./subscribed.json"
)

var poolTopics = sync.Pool{
	New: func() interface{} {
		topics := []Topic{}
		return &topics
	},
}

func init() {
	poolTopics.Put(LoadTopics())
}

// GetTopics get topics from pool
func GetTopics() []Topic {
	topics := poolTopics.Get().(*[]Topic)
	cpTopics := make([]Topic, len(*topics))
	copy(cpTopics, *topics)
	poolTopics.Put(topics)
	return cpTopics
}

// PutTopics id writing file and pool topics reload
func PutTopics(topics []Topic) error {
	byte, err := json.MarshalIndent(topics, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(poolPath, byte, 0666); err != nil {
		return err
	}
	ReLoadTopics()
	return nil
}

// LoadTopics is loading file registered topic information
func LoadTopics() *[]Topic {
	file, err := os.OpenFile(poolPath, os.O_RDWR|os.O_CREATE, 0666)
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

	if len(v) == 0 {
		return &[]Topic{}
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
	poolTopics.Get()
	topics := LoadTopics()
	poolTopics.Put(topics)
}
