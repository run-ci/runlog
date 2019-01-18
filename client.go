package runlog

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
)

// Client is a mediator for communication with Runlog.
type Client struct {
	URL      string
	CertPath string
	KeyPath  string
	CAPath   string

	TaskID uint32

	conn net.Conn
}

// Connect establishes the underlying TCP connection to the Runlog server.
func (c *Client) Connect() error {
	cabuf, err := ioutil.ReadFile(c.CAPath)
	if err != nil {
		return err
	}

	ca := x509.NewCertPool()
	ok := ca.AppendCertsFromPEM(cabuf)
	if !ok {
		return err
	}

	cert, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return err
	}

	tlscfg := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      ca,
	}

	conn, err := tls.Dial("tcp", c.URL, &tlscfg)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) Write(buf []byte) (int, error) {
	packet := Packet{
		TaskID:  c.TaskID,
		Payload: make([]byte, len(buf)),
	}

	bytelen := len(buf)
	if bytelen > 0xFF {
		return -1, errors.New("buffer too large")
	}
	if bytelen < 0 {
		return -1, errors.New("buffer length cannot be negative")
	}

	packet.ByteLength = uint8(bytelen)
	copy(packet.Payload, buf)

	data, err := packet.Encode()
	if err != nil {
		return -1, err
	}

	return c.conn.Write(data)
}
