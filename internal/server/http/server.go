package internalhttp

import (
	"context"
	"net/http"
	"time"

	application "github.com/SashaMelva/calendar_service/internal/app"
)

type Server struct {
	HttpServer *http.Server
	// Addr         string
	// Handler      *Application
	// ReadTimeout  time.Duration
	// WriteTimeout time.Duration
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, app application.App) *Server {
	return &Server{
		&http.Server{
			Addr: app.Host + ":" + app.Port, //add port and host fron config
			// Handler:      app,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.HttpServer.ListenAndServe()
	// <-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
