package main

import (
	"fmt"
	"log"
	"net/http"

	"flag"

	"github.com/midorigreen/gosns/meta"
	"github.com/midorigreen/gosns/topic"
)

func main() {
	fPort := flag.Int("p", 8888, "set server port number")
	flag.Parse()

	if err := run(*fPort); err != nil {
		log.Fatalf("err: %s", err)
	}
}

func run(port int) error {
	topic.Handler()
	meta.Handler()
	err := http.ListenAndServe(":"+fmt.Sprint(port), nil)
	if err != nil {
		return err
	}
	log.Printf("Server running port: %d \n", port)
	return nil
}
