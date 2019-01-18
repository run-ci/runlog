package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net"

	"github.com/run-ci/runlog"
)

// Logger is an interface allowing consumers of this package to
// specify their own logger or just allow the package to log itself.
type Logger interface {
	Printf(format string, v ...interface{})
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
		conn, err := l.Accept()
		if err != nil {
			logger.Printf("got error accepting connection: %v, skipping...")
			continue
		}

		logger.Printf("accepted connection")

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 0xFF)
			packet := runlog.Packet{}

			for {
				n, err := conn.Read(buf)
				if err == io.EOF {
					packet.Decode(buf)
					logger.Printf("GOT EOF %v %+v", n, packet)

					break
				}
				if err != nil {
					logger.Printf("got unknown error: %v", err)

					break
				}

				err = packet.Decode(buf)
				if err == runlog.ErrShortPacket {
					logger.Printf("got short packet %+v", packet)

					break
				}
				if err != nil {
					logger.Printf("got error decoding packet: %v", err)

					break
				}

				logger.Printf("got packet %+v", packet)
			}

		}(conn)
	}
}