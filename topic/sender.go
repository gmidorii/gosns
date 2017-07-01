package topic

import "log"

func sends(topic Topic, topicReq topicReq) {
	for _, v := range topic.Subscribers {
		var service SubscriberService
		switch v.Method.Format {
		case Slack:
			service = &SlackSender{
				URL: v.Method.WebFookURL,
			}
		// case Mail:
		default:
			continue
		}
		if err := service.Send(topicReq.Data); err != nil {
			log.Println(err)
		}
	}
}
