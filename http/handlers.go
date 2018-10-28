package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/hpcloud/tail"
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

func (s *Server) handleGetLog(rw http.ResponseWriter, req *http.Request) {
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
			"task":    task,
			"logsdir": s.logsdir,
		})
		logger.Debug("upgrading request to websocket")

		conn, err := s.upg.Upgrade(rw, req, nil)
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to upgrade to websocket connection")

			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		logpath := fmt.Sprintf("%v/%v.log", s.logsdir, task)

		logger = logger.CloneWith(map[string]interface{}{
			"logpath": logpath,
		})

		logger.Debug("begin streaming task log file")

		t, err := tail.TailFile(logpath, tail.Config{Follow: true})
		if err != nil {
			logger.CloneWith(map[string]interface{}{
				"error": err,
			}).Error("unable to tail file")

			conn.WriteControl(websocket.CloseMessage, []byte(err.Error()), time.Now().Add(2*time.Second))
			return
		}

		for line := range t.Lines {
			l := fmt.Sprintf("%v\n", line.Text)
			err := conn.WriteMessage(websocket.TextMessage, []byte(l))
			if err != nil {
				logger.CloneWith(map[string]interface{}{
					"error": err,
				}).Error("something happened while streaming")

				conn.WriteControl(websocket.CloseMessage, []byte(err.Error()), time.Now().Add(2*time.Second))
				return
			}
		}

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
