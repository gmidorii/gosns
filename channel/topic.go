package channel

// Topic is main struct in package channel
// This has channle name and subscriber list
type Topic struct {
	Channel     string
	Subscribers []Subscriber
}

// Subscriber has subscriber imformation only
type Subscriber struct {
	ClientID string
	Method   Method
}

type Method struct {
	Format     Format
	WebFookURL string
}

// Format is defined as subscribed format
type Format int

const (
	Slack Format = iota
	Mail
	Error
)

// SubscriberService interface must be implemented for subscriber
type SubscriberService interface {
	Send(body string, method Method) error
}

func FormatValue(s string) Format {
	switch s {
	case "slack":
		return Slack
	case "mail":
		return Mail
	default:
		return Error
	}
}
