package main

import (
	"log"
	"net/http"

	"github.com/midorigreen/gosns/channel"
	"github.com/midorigreen/gosns/meta"
)

func main() {
	channel.Handler()
	meta.Handler()
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
	log.Println("Server running port:8888")
}
