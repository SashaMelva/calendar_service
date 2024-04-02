package hendler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	application "github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/storage"
	"go.uber.org/zap"
)

type Service struct {
	Logger zap.SugaredLogger
	app    application.App
	// Ctx    context.Context
	sync.RWMutex
}

type ResponseBody struct {
	Message      string
	MessageError string
}

func NewService(log *zap.SugaredLogger, app *application.App, timeout time.Duration) *Service {
	return &Service{
		Logger: *log,
		app:    *app,
		// Ctx:    ctx,
	}
}

func (s *Service) HendlerEvent(w http.ResponseWriter, req *http.Request) {
	// resp := &ResponseBody{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if req.URL.Path == "/event/" {
		switch req.Method {
		case http.MethodPost:
			s.createEventHandler(w, req, ctx)
		case http.MethodPut:
			s.editEventHandler(w, req)
		case http.MethodGet:
			s.getAllEventsHandler(w, req, ctx)
		case http.MethodDelete:
			s.deleteAllEventsHandler(w, req)
		default:
			s.Logger.Error(fmt.Sprintf("expect method GET, DELETE or POST at /event/, got %v", req.Method))
			return
		}
	} else {
		args := req.URL.Query()
		id := args.Get("id")
		// date := args.Get("date")

		if len(id) > 0 {
			intId, err := strconv.Atoi(id)
			if err != nil {
				s.Logger.Error(fmt.Sprintf("is not valid if event id, got %v", id))
				// resp.Error.Message = fmt.Sprintf("is not valid if event id, got %v", id)
				// w.WriteHeader()
				http.Error(w, err.Error(), http.StatusBadRequest)
				// w.WriteHeader(http.StatusBadRequest)
				// resp.MessageError = fmt.Sprintf("is not valid if event id, got %v", id)
				// js, _ := json.Marshal(resp)
				// w.Write(js)
				return
			}

			switch req.Method {
			case http.MethodDelete:
				s.deleteEventHandlerById(w, req, intId)
			case http.MethodGet:
				s.getEventHandlerById(w, req, intId)
			default:
				s.Logger.Error(fmt.Sprintf("expect method GET or DELETE at /event/<id>, got %v", req.Method))
				return
			}
		}

	}
}

func (s *Service) getAllEventsHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	s.Logger.Info("handling get all events at %s\n", req.URL.Path)

	allEvents, err := s.app.GetAllEvents(ctx)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	js, err := json.Marshal(allEvents)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (s *Service) createEventHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	s.Logger.Info("add new event at %v\n", req.URL.Path)
	event := &storage.Event{}

	err := s.app.CreateEvent(ctx, event)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

}

func (s *Service) editEventHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) deleteAllEventsHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) deleteEventHandlerById(w http.ResponseWriter, req *http.Request, id int) {
}

func (s *Service) getEventHandlerById(w http.ResponseWriter, req *http.Request, id int) {
}
