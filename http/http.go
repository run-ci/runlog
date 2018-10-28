package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/juicemia/log"
)

var logger *log.Logger

var errNoTaskID = errors.New("no task id in path")

func init() {
	logger = log.New("http")

	logger.Debug("logger initialized")
}

type ctxkey int

const (
	keyReqID ctxkey = iota
)

type middleware func(http.HandlerFunc) http.HandlerFunc

func chain(f http.HandlerFunc, mw ...middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		m := mw[i]

		f = m(f)
	}

	return http.HandlerFunc(f)
}

func reqLog(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		reqid := req.Context().Value(keyReqID).(string)

		logger := logger.WithField("request_id", reqid)

		logger.Infof("%v %v", req.Method, req.URL)

		fn(rw, req)
	}
}

func reqID(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		id := uuid.New().String()

		ctx := context.WithValue(req.Context(), keyReqID, id)
		logger.CloneWith(map[string]interface{}{
			"request_id": id,
		}).Debug("setting request ID")

		fn(rw, req.WithContext(ctx))
	}
}

func (s *Server) auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		reqid := req.Context().Value(keyReqID).(string)
		logger := logger.WithField("request_id", reqid)

		user, pass, ok := req.BasicAuth()
		if !ok {
			logger.Error("unauthorized")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if user != s.user {
			logger.Error("unauthorized")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		if pass != s.pass {
			logger.Error("unauthorized")
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		fn(rw, req)
	}
}
