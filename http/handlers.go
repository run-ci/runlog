package http

import (
	"net/http"
	"strings"
)

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

func handleGetLog(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		logger.Debugf("handling GET %v", req.URL)

		task, err := getTaskID(req.URL.String())
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to get task id")
		}

		logger := logger.CloneWith(map[string]interface{}{
			"task": task,
		})
		logger.Debug("getting task log")

		rw.WriteHeader(http.StatusOK)
		return
	default:
		logger.Debugf("can't handle request")

		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func getTaskID(path string) (string, error) {
	segs := strings.Split(path, "/")

	// Split is returning empty strings around the slashes if they are
	// at the beginning or the end of the path.
	if len(segs) >= 3 && segs[2] != "" {
		return segs[2], nil
	}

	return "", errNoTaskID
}
