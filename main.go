package main

import (
	"fmt"
	"net/http"

	"flag"

	"github.com/midorigreen/gosns/logging"
	"github.com/midorigreen/gosns/meta"
	"github.com/midorigreen/gosns/topic"
)

func main() {
	fPort := flag.Int("p", 8888, "set server port number")
	flag.Parse()

	if err := run(*fPort); err != nil {
		logging.Logger.Error(err.Error())
	}
}

func run(port int) error {
	topic.Handler()
	meta.Handler()
	logging.Logger.Info(fmt.Sprintf("Server running port: %d \n", port))

	err := http.ListenAndServe(":"+fmt.Sprint(port), nil)
	if err != nil {
		return err
	}
	return nil
}
