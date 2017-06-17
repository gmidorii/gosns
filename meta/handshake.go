package meta

import (
	"log"
	"math/rand"
	"net/http"
)

type handShakeReq struct {
	Channel                  string
	SupportedConnectionTypes []string
}

type handShakeRes struct {
	Channel                  string
	SupportedConnectionTypes []string
	ClientID                 string `json:"client_id"`
	Successful               bool
}

const (
	randLetter = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	randLen    = 10
)

func handshake(w http.ResponseWriter, r *http.Request) {
	req := handShakeReq{}
	err := decodeBody(r, &req)
	if err != nil {
		log.Println(err)
		res := handShakeRes{
			Channel:    handshakeHandler,
			Successful: false,
		}
		writeRes(res, w)
		return
	}

	res := handShakeRes{
		Channel:                  handshakeHandler,
		SupportedConnectionTypes: req.SupportedConnectionTypes,
		Successful:               true,
		ClientID:                 geneRandStr(randLen),
	}
	writeRes(res, w)
}

func geneRandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = randLetter[rand.Intn(len(randLetter))]
	}
	return string(b)
}
