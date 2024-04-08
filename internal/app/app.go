package app

import (
	"context"
	"fmt"
	"time"

	config "github.com/SashaMelva/calendar_service/internal/config"
	storage "github.com/SashaMelva/calendar_service/internal/storage"
	memorystorage "github.com/SashaMelva/calendar_service/internal/storage/memory"
	"go.uber.org/zap"
)

type Period string

const (
	Day    Period = "Day"
	Week   Period = "Week"
	Mounth Period = "Mounth"
)

type App struct {
	Host    string
	Port    string
	storage *memorystorage.Storage
	Logger  *zap.SugaredLogger
}

type Storage interface {
	GetAllEvents() error
	CreateEvent(storage.Event) error
	DeleteEventById(int) error
	EditEvent(storage.Event) error
	GetByIdEvent(int) error
}

func New(logger *zap.SugaredLogger, storage *memorystorage.Storage, conf *config.ConfigApp) *App {
	return &App{
		Host:    conf.Host,
		Port:    conf.Port,
		storage: storage,
		Logger:  logger,
	}
}
func (a *App) GetAllEvents(ctx context.Context) ([]storage.Event, error) {
	event, err := a.storage.GetAllEvents()

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (a *App) GetByIdEvent(ctx context.Context, id int) (*storage.Event, error) {
	event, err := a.storage.GetByIdEvent(id)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (a *App) CreateEvent(ctx context.Context, event *storage.Event) error {
	id, err := a.storage.CreateEvent(event)

	if err != nil {
		a.Logger.Error(err)
	} else {
		a.Logger.Info(fmt.Sprintf("Create event whith id = %v", id))
	}

	return nil
}
func (a *App) DeleteEventById(id int) error {
	err := a.storage.DeleteEventById(id)

	if err != nil {
		a.Logger.Error(err)
	} else {
		a.Logger.Info(fmt.Sprintf("Delet event whith id = %v", id))
	}

	return err
}
func (a *App) EditEvent(ctx context.Context, event *storage.Event) error {
	err := a.storage.EditEvent(event)

	if err != nil {
		a.Logger.Error(err)
	} else {
		a.Logger.Info(fmt.Sprintf("Update event whith id = %v", event.ID))
	}

	return err
}

func (a *App) GetEventByPeriod(period Period, startDate *time.Time) ([]storage.Event, error) {
	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		a.Logger.Error(err)
	}

	startDateTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, l)
	endDateTime := startDateTime.Add(24 * time.Hour)

	switch period {
	case "Week":
		switch startDateTime.Weekday() {
		case time.Monday:
			endDateTime = startDateTime.Add(7 * 24 * time.Hour)
		case time.Tuesday:
			endDateTime = startDateTime.Add(6 * 24 * time.Hour)
		case time.Wednesday:
			endDateTime = startDateTime.Add(5 * 24 * time.Hour)
		case time.Thursday:
			endDateTime = startDateTime.Add(4 * 24 * time.Hour)
		case time.Friday:
			endDateTime = startDateTime.Add(3 * 24 * time.Hour)
		case time.Saturday:
			endDateTime = startDateTime.Add(2 * 24 * time.Hour)
		case time.Sunday:
			endDateTime = startDateTime.Add(1 * 24 * time.Hour)
		}
	case "Mounth":
		endDateTime = time.Date(startDate.Year(), startDateTime.Month()+1, startDate.Day(), 0, 0, 0, 0, l)
	}

	events, err := a.storage.ListEventsDateForPeriod(&startDateTime, &endDateTime)

	if err != nil {
		a.Logger.Error(err)
	} else {
		a.Logger.Info(fmt.Sprintf("Get events by period %v - %v", startDate, endDateTime))
	}

	return events, nil
}
