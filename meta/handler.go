package meta

import (
	"encoding/json"
	"net/http"
)

var (
	handshakeHandler = "/meta/handshake"
	subHandler       = "/meta/subscribe"
)

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(handshakeHandler, handshake)
	http.HandleFunc(subHandler, subscribe)
}

func writeRes(v interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(v)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
