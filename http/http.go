package http

import (
	"errors"

	"github.com/juicemia/log"
)

var logger *log.Logger
var errNoTaskID = errors.New("no task id in path")

func init() {
	logger = log.New("http")

	logger.Debug("logger initialized")
}
