package http

import "net/http"

func handleRoot(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		logger.Debugf("handling GET %v", req.URL)

		rw.WriteHeader(http.StatusNoContent)
		return
	default:
		logger.Debugf("can't handle request")

		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
