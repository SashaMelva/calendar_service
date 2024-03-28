package hendler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	application "github.com/SashaMelva/calendar_service/internal/app"
	"go.uber.org/zap"
)

type Service struct {
	Logger zap.SugaredLogger
	app    application.App
	sync.RWMutex
}

func NewService(log *zap.SugaredLogger, app *application.App) *Service {
	return &Service{
		Logger: *log,
		app:    *app,
	}
}

func (s *Service) HendlerEvent(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/event/" {
		switch req.Method {
		case http.MethodPost:
			s.createEventHandler(w, req)
		case http.MethodPut:
			s.editEventHandler(w, req)
		case http.MethodGet:
			s.getAllEventsHandler(w, req)
		case http.MethodDelete:
			s.deleteAllEventsHandler(w, req)
		default:
			s.Logger.Error(fmt.Sprintf("expect method GET, DELETE or POST at /event/, got %v", req.Method))
			return
		}
	} else {
		path := strings.Trim(req.URL.Path, "/")
		pathParts := strings.Split(path, "/")

		if len(pathParts) < 2 {
			s.Logger.Error("expect /event/<id> in event handler")
			return
		}

		id, err := strconv.Atoi(pathParts[1])

		if err != nil {
			s.Logger.Error(err)
			return
		}

		switch req.Method {
		case http.MethodDelete:
			s.deleteEventHandler(w, req, int(id))
		case http.MethodGet:
			s.getEventHandler(w, req, int(id))
		default:
			s.Logger.Error(fmt.Sprintf("expect method GET or DELETE at /event/<id>, got %v", req.Method))
			return
		}
	}
}

func (s *Service) getAllEventsHandler(w http.ResponseWriter, req *http.Request) {
	s.Logger.Info("handling get all tasks at %s\n", req.URL.Path)

	allTasks := s.app.GetAllEvents()
	js, err := json.Marshal(allTasks)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (s *Service) createEventHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) editEventHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) deleteAllEventsHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) deleteEventHandler(w http.ResponseWriter, req *http.Request) {
}

func (s *Service) getEventHandler(w http.ResponseWriter, req *http.Request) {
}
