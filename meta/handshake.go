package meta

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

type handShakeRequest struct {
	Channel                  string
	SupportedConnectionTypes []string
}

type handShakeResponse struct {
	Channel                  string
	SupportedConnectionTypes []string
	ClientID                 string `json:"client_id"`
	Successful               bool
}

const (
	randLetter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func handshake(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		res := handShakeResponse{
			Channel:    handshakeHandler,
			Successful: false,
		}
		writeRes(res, w)
		return
	}
	defer r.Body.Close()

	req := handShakeRequest{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
		res := handShakeResponse{
			Channel:    handshakeHandler,
			Successful: false,
		}
		writeRes(res, w)
		return
	}

	res := handShakeResponse{
		Channel:                  handshakeHandler,
		SupportedConnectionTypes: req.SupportedConnectionTypes,
		Successful:               true,
		ClientID:                 randString(10),
	}
	writeRes(res, w)
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randLetter[rand.Intn(len(randLetter))]
	}
	return string(b)
}