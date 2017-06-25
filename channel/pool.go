package channel

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

const (
	poolPath = "./subscribed.json"
)

// PoolTopics is singleton topics cache
var PoolTopics = TopicPool{
	Value: sync.Pool{
		New: func() interface{} {
			topics, err := loadFile(poolPath)
			if err != nil {
				return &[]Topic{}
			}
			return &topics
		},
	},
	Path: poolPath,
}

func init() {
	PoolTopics.Load()
}

// TopicPool is topic info pooling
// topic info master date has `subscribed.json`
type TopicPool struct {
	Value sync.Pool
	Path  string
}

// Pooler is pool interface
type Pooler interface {
	Get() interface{}
	Put(v interface{}) error
	Load()
}

// Get get topics from pool
func (p *TopicPool) Get() interface{} {
	topics := p.Value.Get().(*[]Topic)
	cpTopics := make([]Topic, len(*topics))
	copy(cpTopics, *topics)
	p.Value.Put(topics)
	return cpTopics
}

// Put is writing file and pool topics reload
func (p *TopicPool) Put(v interface{}) error {
	topics, ok := v.([]Topic)
	if ok != true {
		return errors.New("Put argument type []Topic")
	}

	byte, err := json.MarshalIndent(topics, "", "\t")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(p.Path, byte, 0666); err != nil {
		return err
	}
	p.Value.Put(&topics)
	return nil
}

// Load is loading file registered topic information
func (p *TopicPool) Load() error {
	s, err := loadFile(p.Path)
	if err != nil {
		return err
	}
	p.Value.Put(&s)
	return nil
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
