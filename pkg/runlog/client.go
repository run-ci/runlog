package runlog

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type StatusError struct {
	Message    string
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("%v: %v", e.StatusCode, e.Message)
}

type Client struct {
	addr string
	user string
	pass string

	client *http.Client
	upg    *websocket.Upgrader
}

func NewClient(addr, user, pass string) *Client {
	return &Client{
		addr: addr,
		user: user,
		pass: pass,

		client: http.DefaultClient,
	}
}

func (c *Client) GetRoot() error {
	url := fmt.Sprintf("http://%v", c.addr)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return StatusError{
			Message:    err.Error(),
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *Client) GetLog(taskID int, w io.Writer) error {
	url := fmt.Sprintf("ws://%v/log/%v", c.addr, taskID)

	auth := base64.StdEncoding.EncodeToString([]byte(c.user + ":" + c.pass))

	h := http.Header{
		"Authorization": []string{
			fmt.Sprintf("Basic %v", auth),
		},
	}

	conn, resp, err := websocket.DefaultDialer.Dial(url, h)
	if resp == nil {
		if err != nil {
			return err
		}

		return errors.New("got nil response")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return StatusError{
			Message:    "unauthorized",
			StatusCode: resp.StatusCode,
		}
	}

	for {
		typ, msg, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		if typ == websocket.CloseMessage {
			conn.Close()

			break
		}

		_, err = w.Write(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
