package channel

import (
	"sync"
	"testing"
)

func createInitPool() *TopicPool {
	var topicPool = TopicPool{
		Value: sync.Pool{
			New: func() interface{} {
				topics := []Topic{}
				return &topics
			},
		},
		Path: ".././subscribed-sample.json",
	}
	topicPool.Load()
	return &topicPool
}

func TestGet(t *testing.T) {
	topicPool := createInitPool()
	topics := topicPool.Get().([]Topic)

	if len(topics) != 2 {
		t.Error("unexpected topics len :" + string(len(topics)))
	}
	topic0 := topics[0]
	if topic0.Channel != "/news" {
		t.Error("unexpected channel :" + topic0.Channel)
	}
	if len(topic0.Subscribers) != 1 {
		t.Error("unexpected subscribers len :" + string(len(topic0.Subscribers)))
	}
	subscriver0 := topic0.Subscribers[0]
	if subscriver0.ClientID != "XVlBzgbaiC" {
		t.Error("unexpected clientID :" + subscriver0.ClientID)
	}
	if subscriver0.Method.Format != 0 {
		t.Error("unexpected method format :" + string(subscriver0.Method.Format))
	}
	if subscriver0.Method.WebFookURL != "http://hoge" {
		t.Error("unexpected method webfookurl :" + subscriver0.Method.WebFookURL)
	}
}

func TestGetLoop(t *testing.T) {
	topicPool := createInitPool()
	for i := 0; i < 3; i++ {
		topics := topicPool.Get().([]Topic)

		if len(topics) != 2 {
			t.Error("unexpected topics len :" + string(len(topics)))
		}
		topic0 := topics[0]
		if topic0.Channel != "/news" {
			t.Error("unexpected channel :" + topic0.Channel)
		}
		if len(topic0.Subscribers) != 1 {
			t.Error("unexpected subscribers len :" + string(len(topic0.Subscribers)))
		}
		subscriver0 := topic0.Subscribers[0]
		if subscriver0.ClientID != "XVlBzgbaiC" {
			t.Error("unexpected clientID :" + subscriver0.ClientID)
		}
		if subscriver0.Method.Format != 0 {
			t.Error("unexpected method format :" + string(subscriver0.Method.Format))
		}
		if subscriver0.Method.WebFookURL != "http://hoge" {
			t.Error("unexpected method webfookurl :" + subscriver0.Method.WebFookURL)
		}
	}
}
