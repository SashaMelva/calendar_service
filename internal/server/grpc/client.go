package internalgrpc

import (
	"github.com/SashaMelva/calendar_service/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	proto "github.com/SashaMelva/calendar_service/internal/server/grpc/gen"
)

func NewGRPCClient(config *config.ConfigGrpcServer, log *zap.SugaredLogger) proto.EventServiceClient {
	conn, err := grpc.Dial(config.Host + ":" + config.Port)

	if err != nil {
		log.Error(err)
	}

	client := proto.NewEventServiceClient(conn)

	return client
}
