package main

import (
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/meta"
	"github.com/midorigreen/gopubsub/topic"
)

func main() {
	topic.Handler()
	meta.Handler()
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
