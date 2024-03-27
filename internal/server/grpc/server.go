package internalgrpc

import (
	"context"
	"time"
)

type Server struct {
	Addr         string
	Handler      *Application
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	return &Server{
		Addr:         ":", //add port and host fron config
		Handler:      &app,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}
