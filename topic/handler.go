package topic

import (
	"log"
	"net/http"
)

const (
	topicsPath = "/topics"
)

// Handler is topic handler
func Handler() {
	http.HandleFunc(topicsPath, topicsHandler)
}

func topicsHandler(w http.ResponseWriter, r *http.Request) {
	m := Mail{}
	err := m.Send("")
	if err != nil {
		log.Fatal(err)
	}
}
