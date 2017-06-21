package channel

import (
	"net/smtp"
)

type Mail struct {
	from    string
	to      string
	title   string
	body    string
	usrname string
	pwd     string
}

func (m *Mail) Send(body string) error {
	auth := smtp.PlainAuth("", m.usrname, m.pwd, "smtp.gmail.com")
	if err := smtp.SendMail("smtp.gmail.com:587", auth, m.from, []string{m.to}, []byte(m.body)); err != nil {
		return err
	}
	return nil
}
