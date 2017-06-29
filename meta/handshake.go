package meta

import (
	"log"
	"math/rand"
	"net/http"
)

type handShakeReq struct {
	Channel string `json:"channel"`
}

type handShakeRes struct {
	Channel    string `json:"channel"`
	ClientID   string `json:"client_id"`
	Successful bool   `json:"successful"`
}

const (
	randLetter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randLen    = 10
)

func handshakeHandler(w http.ResponseWriter, r *http.Request) {
	req := handShakeReq{}
	err := decodeBody(r, &req)
	if err != nil {
		log.Println(err)
		res := handShakeRes{
			Channel:    handshakePattarn,
			Successful: false,
		}
		writeRes(res, w, http.StatusBadRequest)
		return
	}

	res := handShakeRes{
		Channel:    handshakePattarn,
		Successful: true,
		ClientID:   geneRandStr(randLen),
	}
	writeRes(res, w, http.StatusOK)
}

func geneRandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randLetter[rand.Intn(len(randLetter))]
	}
	return string(b)
}
