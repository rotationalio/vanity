package server

import "github.com/rotationalio/vanity/config"

type Server struct{}

func New(conf config.Config) (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Serve() error {
	return nil
}
