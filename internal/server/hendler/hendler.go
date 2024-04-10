package hendler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	application "github.com/SashaMelva/calendar_service/internal/app"
	"github.com/SashaMelva/calendar_service/internal/server/validate"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if req.Method == http.MethodGet {

		if req.URL.Path == "/event/day" {
			args := req.URL.Query()
			date := args.Get("date")
			s.getEventsByPeriod(w, date, "Day")
		} else if req.URL.Path == "/event/week" {
			args := req.URL.Query()
			date := args.Get("date")
			s.getEventsByPeriod(w, date, "Week")
		} else if req.URL.Path == "/event/mounth" {
			args := req.URL.Query()
			date := args.Get("date")
			s.getEventsByPeriod(w, date, "Mounth")
		}
	}

	if req.URL.Path == "/event/" {
		switch req.Method {
		case http.MethodGet:
			s.getAllEventsHandler(w, req, ctx)
		case http.MethodPost:
			s.createEventHandler(w, req, ctx)
		case http.MethodPut:
			s.editEventHandler(w, req, ctx)
		default:
			s.Logger.Error(fmt.Sprintf("expect method GET, DELETE or POST at /event/, got %v", req.Method))
			return
		}
	} else {
		args := req.URL.Query()
		id := args.Get("id")
		// date := args.Get("date")

		if len(id) > 0 {
			s.Logger.Info("id event " + id)
			intId, err := strconv.Atoi(id)
			if err != nil {
				s.Logger.Error(fmt.Sprintf("is not valid if event id, got %v", id))
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			switch req.Method {
			case http.MethodDelete:
				s.deleteEventHandlerById(w, intId)
			case http.MethodGet:
				s.getEventHandlerById(w, ctx, intId)
			default:
				s.Logger.Error(fmt.Sprintf("expect method GET or DELETE at /event?=<id>, got %v", req.Method))
				return
			}
		}

	}
}

func (s *Service) getAllEventsHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	s.Logger.Info("handling get all events at %s\n", req.URL.Path)

	allEvents, err := s.app.GetAllEvents(ctx)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(allEvents)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (s *Service) getEventsByPeriod(w http.ResponseWriter, strDate, period string) {
	s.Logger.Info("handling get events by day %s\n")

	date, err := time.Parse(strDate, "")

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var p application.Period
	p = "Day"

	switch period {
	case "Week":
		p = "Week"
	case "Mounth":
		p = "Mounth"
	}

	allEvents, err := s.app.GetEventByPeriod(p, &date)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(allEvents)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}

func (s *Service) createEventHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	s.Logger.Info("add new event at %v\n", req.URL.Path)

	event := storage.Event{}
	body, err := io.ReadAll(req.Body)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		err = json.Unmarshal(body, &event)
		if err != nil {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
		}
	}

	ok, msg := validate.ValidEvent(&event)

	if ok != "OK" {
		s.Logger.Error("dont valid data")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
		return
	}

	errCreater := s.app.CreateEvent(ctx, &event)

	if errCreater != nil {
		s.Logger.Error(w, errCreater.Error(), http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("dont created"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) editEventHandler(w http.ResponseWriter, req *http.Request, ctx context.Context) {
	s.Logger.Info("edit event at %v\n", req.URL.Path)
	fmt.Println("qweqewqe")

	event := storage.Event{}
	body, err := io.ReadAll(req.Body)

	if err != nil {
		s.Logger.Error(w, "err %q\n", err.Error())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
	} else {
		err = json.Unmarshal(body, &event)
		if err != nil {
			s.Logger.Error(w, "can't unmarshal: ", err.Error())
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("can't unmarshal: " + err.Error()))
		}
	}

	if event.ID == 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Id event empty"))
		return
	}

	ok, msg := validate.ValidEvent(&event)

	if ok != "OK" {
		s.Logger.Error("dont valid data")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(msg))
		return
	}

	s.Logger.Info(event)

	errEdit := s.app.EditEvent(ctx, &event)

	if errEdit != nil {
		s.Logger.Error(w, errEdit.Error(), http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("dont update"))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) deleteEventHandlerById(w http.ResponseWriter, id int) {
	s.Logger.Info("delet event by id %v", id)

	err := s.app.DeleteEventById(id)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) getEventHandlerById(w http.ResponseWriter, ctx context.Context, id int) {
	s.Logger.Info("handling get event at by id %v", id)

	event, err := s.app.GetByIdEvent(ctx, id)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	js, err := json.Marshal(event)

	if err != nil {
		s.Logger.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
}
