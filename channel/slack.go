package channel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SlackSender is implemented SubscriberService interface
type SlackSender struct {
	URL string
}

type slackBody struct {
	Text string `json:"text"`
}

const bodyType = "application/json"

// Send has method send message to slack
func (s *SlackSender) Send(body string) error {
	fmt.Println(body)
	sb := slackBody{
		Text: body,
	}
	postBody, err := json.Marshal(sb)
	if err != nil {
		return err
	}
	resp, err := http.Post(s.URL, bodyType, bytes.NewBuffer(postBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("bad request status code: " + string(resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if string(respBody) != "ok" {
		return errors.New("failed request sending")
	}
	return nil
}
