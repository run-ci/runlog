package http

import (
	"errors"
	"net/http"
)

// Server is an HTTP server that can serve requests for logs.
type Server struct {
	user string
	pass string

	*http.Server
}

// NewServer initializes and returns a *Server. The server is only
// accessible using basic auth with the credentials passed as
// `user` and `pass`.
func NewServer(user, pass string) (*Server, error) {
	if user == "" {
		return nil, errors.New("must specify a user")
	}

	if pass == "" {
		return nil, errors.New("must specify a password")
	}

	mux := http.NewServeMux()

	mux.Handle("/", http.HandlerFunc(handleRoot))

	return &Server{
		user: user,
		pass: pass,

		Server: &http.Server{
			Handler: mux,
			Addr:    ":7777",
		},
	}, nil
}
