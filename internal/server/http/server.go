package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/config"
	hendler "github.com/SashaMelva/calendar_service/internal/server/hendler"
	"go.uber.org/zap"
)

type Server struct {
	HttpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, app *app.App, config *config.ConfigHttpServer) *Server {
	log.Info("URL api" + config.Host + ":" + config.Port)
	timeout := config.Timeout * time.Second

	mux := http.NewServeMux()
	h := hendler.NewService(log, app, timeout)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	mux.HandleFunc("/event/", h.HendlerEvent)

	return &Server{
		&http.Server{
			Addr:         config.Host + ":" + config.Port,
			Handler:      mux,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	err := s.HttpServer.ListenAndServe()
	<-ctx.Done()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	return s.HttpServer.Shutdown(ctx)
}
