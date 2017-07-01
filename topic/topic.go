package topic

// Topic is main struct in package topic
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

// Method is defined as message transmission method
type Method struct {
	Format     Format
	WebFookURL string
}

// Format is defined as subscribed format
type Format int

const (
	// Slack is format type slack service
	Slack Format = iota
	// Mail is format type e-mail
	Mail
	// Error is case error
	Error
)

// SubscriberService interface must be implemented for subscriber
type SubscriberService interface {
	Send(body string) error
}

// FormatValue is string to Format type
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
