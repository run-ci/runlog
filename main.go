package main

import (
	"os"

	"github.com/juicemia/log"
	"github.com/run-ci/runlog/http"
)

var quser, qpass, logsdir string
var logger *log.Logger

func init() {
	log.SetLevelFromEnv("RUNLOG_LOG_LEVEL")
	logger = log.New("main")

	logger.Info("reading environment")

	quser = os.Getenv("RUNLOG_HTTP_USER")
	if quser == "" {
		quser = "runlog_devel"
		logger.Debugf("RUNLOG_HTTP_USER empty, defaulting to %v", quser)
	}

	qpass = os.Getenv("RUNLOG_HTTP_PASS")
	if qpass == "" {
		qpass = "runlog_devel"
		logger.Debugf("RUNLOG_HTTP_PASS empty, defaulting to %v", qpass)
	}

	logsdir = os.Getenv("RUNLOG_LOGS_DIR")
	if logsdir == "" {
		logger.Fatal("must set RUNLOG_LOGS_DIR")
	}
}

func main() {
	logger.Info("booting runlog query service")

	srv, err := http.NewServer(quser, qpass, logsdir)
	if err != nil {
		logger.CloneWith(map[string]interface{}{"error": err}).
			Fatal("unable to start server")
	}

	logger.Fatal(srv.ListenAndServe())
}
