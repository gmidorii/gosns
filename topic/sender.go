package topic

import "github.com/midorigreen/gosns/logging"

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
			logging.Logger.Error(err.Error())
		}
	}
}
