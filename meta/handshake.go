package meta

import (
	"net/http"
)

type request struct {
	Channel                  string
	SupportedConnectionTypes []string
}

type response struct {
	Channel                  string
	SupportedConnectionTypes []string
	ClientID                 string `json:"client_id"`
	Successfult              bool
}

func handshake(w http.ResponseWriter, r *http.Request) {

}
