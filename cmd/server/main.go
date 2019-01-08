package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"

	"github.com/run-ci/runlog"
	"github.com/sirupsen/logrus"
)

var quser, qpass, logsdir, capath, certpath, keypath string
var logger *logrus.Entry

func init() {
	logger = logrus.WithField("package", "main")

	logger.Info("reading environment")

	logsdir = os.Getenv("RUNLOG_LOGS_DIR")
	if logsdir == "" {
		logger.Fatal("must set RUNLOG_LOGS_DIR")
	}

	capath = os.Getenv("RUNLOG_CAPATH")
	if capath == "" {
		logger.Warn("RUNLOG_CAPATH not set, using /etc/runlog/ssl/ca.pem")
		capath = "/etc/runlog/ssl/ca.pem"
	}

	certpath = os.Getenv("RUNLOG_CERTPATH")
	if certpath == "" {
		logger.Warn("RUNLOG_CERTPATH not set, using /etc/runlog/ssl/runlog.crt")
		certpath = "/etc/runlog/ssl/runlog.crt"
	}

	keypath = os.Getenv("RUNLOG_KEYPATH")
	if keypath == "" {
		logger.Warn("RUNLOG_KEYPATH not set, using /etc/runlog/ssl/runlog.key")
		keypath = "/etc/runlog/ssl/runlog.key"
	}
}

func main() {
	logger.Info("booting runlog query service")

	//srv, err := http.NewServer(quser, qpass, logsdir)
	//if err != nil {
	//logger.CloneWith(map[string]interface{}{"error": err}).
	//Fatal("unable to start server")
	//}

	srv := runlog.Server{
		Addr: ":9999",
	}

	cabuf, err := ioutil.ReadFile(capath)
	if err != nil {
		logger.WithError(err).Fatalf("unable to open %v", capath)
	}

	srv.CertificateAuthority = x509.NewCertPool()
	ok := srv.CertificateAuthority.AppendCertsFromPEM(cabuf)
	if !ok {
		logger.Fatalf("unable to append PEM data to CA cert pool")
	}

	srv.Certificate, err = tls.LoadX509KeyPair(certpath, keypath)
	if err != nil {
		logger.WithError(err).Fatal("unable to load X509 key pair")
	}

	logger.Fatal(srv.ListenAndServeTLS())
}
