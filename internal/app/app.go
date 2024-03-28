package app

import (
	"context"

	config "github.com/SashaMelva/calendar_service/internal/config"
)

type App struct {
	Host         string
	Port         string
	ConnectionDB *Storage
}

type Logger interface { // TODO
}

type Storage interface {
}

func New(logger Logger, storage Storage, conf *config.ConfigApp) *App {
	return &App{
		Host:         conf.Host,
		Port:         conf.Port,
		ConnectionDB: &storage,
	}
}
func (a *App) GetAllEvents(ctx context.Context) error {
	// TODO

	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO

	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
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
