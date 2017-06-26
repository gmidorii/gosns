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
	topics, err := loadFile(file)
	if err != nil {
		log.Fatalf("err: %s", err)
	}
	for _, v := range topics {
		subByte, err := json.Marshal(v.Subscribers)
		if err != nil {
			log.Fatalln(err)
		}
		if err = ds.Set(v.Channel, string(subByte), 24*time.Hour); err != nil {
			log.Fatalln(err)
		}
	}
}

// TopicDataService is interface dealing with topic data
type TopicDataService interface {
	Add(topic Topic) error
	Update(topic []Topic) error
	Fetch(channel string) (Topic, error)
}

// TopicData is topic file path and cache
type TopicData struct {
	Path string
	Ds   *dscache.Dscache
}

// CreateTopicData is generating topic data
// TopicData struct should be only created by this func
func CreateTopicData() *TopicData {
	return &TopicData{
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
func (d *TopicData) Update(topics []Topic) error {
	nowTopics, err := loadFile(d.Path)
	if err != nil {
		return err
	}
	var upTopics []Topic
	for _, topic := range topics {
		for i, nowTopic := range nowTopics {
			if topic.Channel == nowTopic.Channel {
				// TODO: duplication check
				nowTopics[i].Subscribers = append(nowTopic.Subscribers, topic.Subscribers...)
				upTopics = append(upTopics, topic)
				break
			}
		}
	}
	if len(upTopics) == 0 {
		return errors.New("not existing topic which try to update")
	}

	// update file
	if err = writeFile(nowTopics, d.Path); err != nil {
		return err
	}

	// update cache
	for _, v := range upTopics {
		// cache clear
		ok := d.Ds.Purge(v.Channel)
		if ok != true {
			return errors.New("failed update cache")
		}
		subByte, err := json.Marshal(v.Subscribers)
		if err != nil {
			return err
		}
		if err = d.Ds.Set(v.Channel, string(subByte), 24*time.Hour); err != nil {
			return err
		}
	}
	return nil
}

// Fetch is fetching topic from cache (and file)
func (t *TopicData) Fetch(channel string) (Topic, error) {
	str, ok := t.Ds.Get(channel)
	if ok != true {
		// TODO: when not found cache, read file
		return Topic{}, errors.New("not found channel from cache")
	}
	var subscribers []Subscriber
	err := json.Unmarshal([]byte(str), &subscribers)
	if err != nil {
		return Topic{}, err
	}
	return Topic{
		Channel:     channel,
		Subscribers: subscribers,
	}, nil
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
