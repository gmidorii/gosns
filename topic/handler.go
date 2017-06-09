package topic

import "net/http"

const topicsPath = "/topics"

// Handler is topic handler
func Handler(path string) {
	http.HandleFunc(path+topicsPath, topics)
}

func topics(w http.ResponseWriter, r *http.Request) {
}
