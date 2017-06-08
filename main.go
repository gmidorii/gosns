package main

import (
	"log"
	"net/http"

	"github.com/midorigreen/gopubsub/meta"
)

func main() {
	meta.Handler()
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}
