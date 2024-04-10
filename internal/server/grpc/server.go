package internalgrpc

import (
	"context"
	"errors"
	"net"

	"github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/config"
	"github.com/SashaMelva/calendar_service/internal/server/validate"
	"github.com/SashaMelva/calendar_service/internal/storage"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	proto "github.com/SashaMelva/calendar_service/internal/server/grpc/gen"
)

type Service struct {
	Logger *zap.SugaredLogger
	app    *app.App
	proto.UnimplementedEventServiceServer
}

func NewGRPCServer(log *zap.SugaredLogger, app *app.App) *grpc.Server {

	gsrv := grpc.NewServer()
	srv := &Service{
		Logger: log,
		app:    app,
	}

	proto.RegisterEventServiceServer(gsrv, srv)
	return gsrv
}

func ListenServer(server *grpc.Server, config *config.ConfigGrpcServer, log *zap.SugaredLogger) {
	log.Info("Starting listening on port " + config.Port)
	port := config.Host + ":" + config.Port

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Info("Listening on %s", port)

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Service) GetEventsByPeriod(ctx context.Context, req *proto.EventDateForPeriodRequest) (*proto.GetEventsResponse, error) {
	date := req.Date.AsTime()
	period := string(req.Period)

	var p app.Period
	p = "Day"

	switch period {
	case "WEEK":
		p = "Week"
	case "MOUNTH":
		p = "Mounth"
	}

	eventsStorage, err := s.app.GetEventByPeriod(p, &date)

	if err != nil {
		return &proto.GetEventsResponse{Event: nil}, err
	}

	eventsPb := make([]*proto.Event, len(eventsStorage))

	for i := range eventsPb {
		event := eventsStorage[i]
		eventsPb[i] = &proto.Event{
			Id:            uint32(event.ID),
			Title:         event.Title,
			DateTimeStart: timestamppb.New(*event.DateTimeStart),
			DateTimeEnd:   timestamppb.New(*event.DateTimeEnd),
			Description:   event.Description,
		}
	}

	return &proto.GetEventsResponse{Event: eventsPb}, nil
}

func (s *Service) GetEventById(ctx context.Context, req *proto.EventId) (*proto.GetEventResponse, error) {
	if req.Id == 0 {
		s.Logger.Error("Id event empty")
		return &proto.GetEventResponse{}, errors.New("Id event empty")
	}

	event, err := s.app.GetByIdEvent(ctx, int(req.Id))

	if err != nil {
		s.Logger.Error(err)
		return &proto.GetEventResponse{}, err
	}

	var timePdStart *timestamppb.Timestamp
	if event.DateTimeStart == nil {
		timePdStart = nil
	} else {
		timePdStart = timestamppb.New(*event.DateTimeStart)
	}

	var timePdEnd *timestamppb.Timestamp
	if event.DateTimeEnd == nil {
		timePdEnd = nil
	} else {
		timePdEnd = timestamppb.New(*event.DateTimeEnd)
	}

	return &proto.GetEventResponse{
		Event: &proto.Event{
			Id:            uint32(event.ID),
			Title:         event.Title,
			DateTimeStart: timePdStart,
			DateTimeEnd:   timePdEnd,
			Description:   event.Description,
		},
	}, nil
}

func (s *Service) CreateEvent(ctx context.Context, req *proto.Event) (*proto.GetResponse, error) {
	dateStart := req.DateTimeStart.AsTime()
	dateEnd := req.DateTimeEnd.AsTime()
	event := storage.Event{
		Title:         req.Title,
		DateTimeStart: &dateStart,
		DateTimeEnd:   &dateEnd,
		Description:   req.Description,
	}
	ok, msg := validate.ValidEvent(&event)

	if ok != "OK" {
		s.Logger.Error("dont valid data")
		return &proto.GetResponse{Status: proto.Status_ERROR}, errors.New(msg)
	}

	err := s.app.CreateEvent(ctx, &event)

	if err != nil {
		return &proto.GetResponse{Status: proto.Status_ERROR}, err
	}

	return &proto.GetResponse{Status: proto.Status_OK}, nil
}

func (s *Service) EditEvent(ctx context.Context, req *proto.Event) (*proto.GetResponse, error) {
	dateStart := req.DateTimeStart.AsTime()
	dateEnd := req.DateTimeEnd.AsTime()
	event := storage.Event{
		ID:            int(req.Id),
		Title:         req.Title,
		DateTimeStart: &dateStart,
		DateTimeEnd:   &dateEnd,
		Description:   req.Description,
	}

	if event.ID == 0 {
		s.Logger.Error("Id event empty")
		return &proto.GetResponse{Status: proto.Status_ERROR}, errors.New("Id event empty")
	}
	ok, msg := validate.ValidEvent(&event)

	if ok != "OK" {
		s.Logger.Error("dont valid data: " + msg)
		return &proto.GetResponse{Status: proto.Status_ERROR}, errors.New(msg)
	}
	err := s.app.CreateEvent(ctx, &event)

	if err != nil {
		return &proto.GetResponse{Status: proto.Status_ERROR}, err
	}

	return &proto.GetResponse{Status: proto.Status_OK}, nil
}

func (s *Service) DeleteEventById(ctx context.Context, req *proto.EventId) (*proto.GetResponse, error) {
	if req.Id == 0 {
		s.Logger.Error("Id event empty")
		return &proto.GetResponse{Status: proto.Status_ERROR}, errors.New("Id event empty")
	}

	err := s.app.DeleteEventById(int(req.Id))

	if err != nil {
		return &proto.GetResponse{Status: proto.Status_ERROR}, err
	}
	return &proto.GetResponse{Status: proto.Status_OK}, nil
}
