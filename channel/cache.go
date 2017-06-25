package channel

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"time"

	"errors"

	"github.com/emluque/dscache"
)

var file = "./subscribed.json"
var ds, dsErr = dscache.New(2 * dscache.MB)

func init() {
	if dsErr != nil {
		log.Fatalf("err: %s", dsErr)
	}
}

// TopicDataService is interface dealing with topic data
type TopicDataService interface {
	Add(topic Topic) error
	Update(topic Topic) error
}

// TopicData is topic file path and cache
type TopicData struct {
	Path string
	Ds   *dscache.Dscache
}

// CreateTopicData is generating topic data
// TopicData struct should be only created by this func
func CreateTopicData() TopicData {
	return TopicData{
		Path: file,
		Ds:   ds,
	}
}

// Add is add topic data to file and cache
func (d *TopicData) Add(topic Topic) error {
	topics, err := loadFile(d.Path)
	if err != nil {
		return err
	}
	topics = append(topics, topic)

	// add file
	if err = writeFile(topics, d.Path); err != nil {
		return err
	}

	// add cache
	subByte, err := json.Marshal(topic.Subscribers)
	if err != nil {
		return err
	}
	return d.Ds.Set(topic.Channel, string(subByte), 24*time.Hour)
}

// Update is update topic to file and cache
func (d *TopicData) Update(topic Topic) error {
	topics, err := loadFile(d.Path)
	if err != nil {
		return err
	}
	var upTopic Topic
	for _, v := range topics {
		if v.Channel == topic.Channel {
			// TODO: duplication check
			v.Subscribers = append(v.Subscribers, topic.Subscribers...)
			upTopic = v
			break
		}
	}
	if upTopic.Channel == "" {
		return errors.New("not existing topic which try to update")
	}

	// update file
	if err = writeFile(topics, d.Path); err != nil {
		return err
	}

	// update cache
	_, ok := d.Ds.Get(upTopic.Channel)
	if ok != true {
		return errors.New("failed update cache")
	}
	subByte, err := json.Marshal(upTopic.Subscribers)
	if err != nil {
		return err
	}
	return d.Ds.Set(upTopic.Channel, string(subByte), 24*time.Hour)
}

func loadFile(path string) ([]Topic, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	v, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	s := []Topic{}
	err = json.Unmarshal(v, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func writeFile(topics []Topic, path string) error {
	byte, err := json.MarshalIndent(topics, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(path, byte, 0666); err != nil {
		return err
	}
	return nil
}
