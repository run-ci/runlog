package runlog

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
)

// Logger is an interface allowing consumers of this package to
// specify their own logger or just allow the package to log itself.
type Logger interface {
	Printf(format string, v ...interface{})
}

var logger Logger

func init() {
	logger = log.New(os.Stderr, "[RUNLOG DEBUG] ", log.LstdFlags)
}

// SetLogger sets the package logger.
func SetLogger(l Logger) {
	logger = l
}

// Server serves logs over TLS/TCP as well as ingesting them
// over TLS/TCP. It holds the information about the certificates
// it needs to communicate with a client.
type Server struct {
	Addr string

	CertificateAuthority *x509.CertPool
	Certificate          tls.Certificate
}

// ListenAndServeTLS blocks serving logging requests and
// ingesting logs over TLS/TCP.
func (srv *Server) ListenAndServeTLS() error {
	logger.Printf("setting TLS configuration")

	cfg := &tls.Config{
		Certificates: []tls.Certificate{
			srv.Certificate,
		},

		RootCAs: srv.CertificateAuthority,
	}

	l, err := tls.Listen("tcp", srv.Addr, cfg)
	if err != nil {
		return err
	}

	logger.Printf("listening for connections on %v", l.Addr())
	for {
		_, err := l.Accept()
		if err != nil {
			logger.Printf("got error accepting connection: %v, skipping...")
			continue
		}

		logger.Printf("accepted connection")
	}
}
