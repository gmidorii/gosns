package channel

import (
	"encoding/json"
	"net/http"
)

const (
	topicsPath = "/topics"
)

type topicReq struct {
	Channel string `json:"channel"`
	Data    string `json:"data"`
}

// Handler is topic handler
func Handler() {
	http.HandleFunc(topicsPath, topicsHandler)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	var t topicReq
	decodeBody(r, &t)
	w.Write([]byte(t.Channel))
}

func decodeBody(req *http.Request, out interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(out)
}
