package http

import (
	"github.com/juicemia/log"
)

var logger *log.Logger

func init() {
	logger = log.New("http")

	logger.Debug("logger initialized")
}
