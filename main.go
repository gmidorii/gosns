package main

import (
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/meta"
	"github.com/midorigreen/gopubsub/topic"
)

const (
	path = "/gopubsub"
)

func main() {
	topic.Handler(path)
	meta.Handler(path)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
