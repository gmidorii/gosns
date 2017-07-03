package logging

import (
	"log"

	"go.uber.org/zap"
)

// Logger is used all code
var Logger, err = zap.NewProduction()

func init() {
	if err != nil {
		log.Fatalf("err: %s", err)
	}
}
