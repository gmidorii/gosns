package main

import (
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/channel"
	"github.com/midorigreen/gopubsub/meta"
)

func main() {
	channel.Handler()
	meta.Handler()
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
