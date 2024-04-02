package app

import (
	"context"
	"fmt"

	config "github.com/SashaMelva/calendar_service/internal/config"
	storage "github.com/SashaMelva/calendar_service/internal/storage"
	memorystorage "github.com/SashaMelva/calendar_service/internal/storage/memory"
	"go.uber.org/zap"
)

type App struct {
	Host    string
	Port    string
	storage *memorystorage.Storage
	Logger  *zap.SugaredLogger
}

type Logger interface { // TODO
}

type Storage interface {
	GetAllEvents() error
	CreateEvent(storage.Event) error
	DeleteEvent(int) error
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
	a.Logger.Info("qwqwqw")
	id, err := a.storage.CreateEvent(event)

	if err != nil {
		a.Logger.Error(err)
	}

	a.Logger.Info(fmt.Sprintf("Create event whith id = %v", id))
	return nil
}
func (a *App) DeleteEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
func (a *App) EditEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
