package meta

import (
	"log"
	"net/http"
)

var (
	subHandler = "/meta/subscribe"
)

// Handler is meta channel definition
func Handler() {
	http.HandleFunc(subHandler, subscribe)
}

func subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println(subHandler)
}
