package internalgrpc

import (
	"context"

	application "github.com/SashaMelva/calendar_service/internal/app"
	"go.uber.org/zap"
)

type Server struct {
}

func NewServer(log *zap.SugaredLogger, app *application.App) *Server {
}

func (s *Server) Start(ctx context.Context) error {
}

func (s *Server) Stop(ctx context.Context) error {
}
