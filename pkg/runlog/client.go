package runlog

import (
	"fmt"
	"net/http"
)

type StatusError struct {
	Message    string
	StatusCode int
}

func (e StatusError) Error() string {
	return fmt.Sprintf("got wrong status code %v", e.StatusCode)
}

type Client struct {
	addr string

	client *http.Client
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,

		client: http.DefaultClient,
	}
}

func (c Client) GetRoot() error {
	req, err := http.NewRequest(http.MethodGet, c.addr, nil)
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
