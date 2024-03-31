package app

import (
	"context"
	"time"

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
	CreateEvent(storage.Event) error
	DeleteEvent(int) error
	EditEvent(storage.Event) error
}

func New(logger *zap.SugaredLogger, storage *memorystorage.Storage, conf *config.ConfigApp) *App {
	return &App{
		Host:    conf.Host,
		Port:    conf.Port,
		storage: storage,
		Logger:  logger,
	}
}
func (a *App) GetAllEvents(ctx context.Context) error {
	// TODO

	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) CreateEvent(ctx context.Context, id, title, descriptioin string, time time.Time, date time.Duration) error {
	// TODO

	return a.storage.CreateEvent(&storage.Event{
		ID:           id,
		Title:        title,
		Date:         date,
		Time:         time,
		Descriptioin: descriptioin,
	})
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
