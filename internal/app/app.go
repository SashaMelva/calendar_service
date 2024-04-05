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
