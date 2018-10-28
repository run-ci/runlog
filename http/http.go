package http

import (
	"errors"
	"net/http"

	"github.com/juicemia/log"
)

var logger *log.Logger

var errNoTaskID = errors.New("no task id in path")

func init() {
	logger = log.New("http")

	logger.Debug("logger initialized")
}

type middleware func(http.HandlerFunc) http.HandlerFunc

func chain(f http.HandlerFunc, mw ...middleware) http.Handler {
	for i := 0; i < len(mw); i++ {
		m := mw[i]

		f = m(f)
	}

	return http.HandlerFunc(f)
}

func reqLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		logger.Infof("%v %v", req.Method, req.URL)

		fn(rw, req)
	}
}
