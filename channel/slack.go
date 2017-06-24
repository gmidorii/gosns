package channel

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

const bodyType = "application/json"

// Send is implemented SubscriberService interface
// this send message to slack
func Send(body string, method Method) error {
	resp, err := http.Post(method.WebFookURL, bodyType, bytes.NewBufferString(body))
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
