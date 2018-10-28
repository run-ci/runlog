package http

import (
	"errors"
	"net/http"
)

// Server is an HTTP server that can serve requests for logs.
type Server struct {
	user    string
	pass    string
	logsdir string

	*http.Server
}

// NewServer initializes and returns a *Server. The server is only
// accessible using basic auth with the credentials passed as
// `user` and `pass`.
func NewServer(user, pass, logsdir string) (*Server, error) {
	if user == "" {
		return nil, errors.New("must specify a user")
	}

	if pass == "" {
		return nil, errors.New("must specify a password")
	}

	if logsdir == "" {
		return nil, errors.New("must specify a logsdir")
	}

	mux := http.NewServeMux()
	srv := &Server{
		user:    user,
		pass:    pass,
		logsdir: logsdir,

		Server: &http.Server{
			Handler: mux,
			Addr:    ":7777",
		},
	}

	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/log/", http.HandlerFunc(srv.handleGetLog))

	return srv, nil
}
