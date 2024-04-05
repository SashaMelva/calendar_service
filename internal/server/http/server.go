package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	application "github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/server/hendler"
	"go.uber.org/zap"
)

type Server struct {
	HttpServer *http.Server
}

type Logger interface { // TODO
}

// type Application interface { // TODO
// }

func NewServer(log *zap.SugaredLogger, app *application.App) *Server {
	log.Info("URL api" + app.Host + ":" + app.Port)
	timeout := 10 * time.Second

	mux := http.NewServeMux()
	h := hendler.NewService(log, app, timeout)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	mux.HandleFunc("/event/", h.HendlerEvent)

	return &Server{
		&http.Server{
			Addr:         app.Host + ":" + app.Port,
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
