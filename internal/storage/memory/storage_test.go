package memorystorage

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	storage "github.com/SashaMelva/calendar_service/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
)

func TestCreateEvent(t *testing.T) {
	testCases := []struct {
		name  string
		event *storage.Event
	}{
		{
			name: "all required fields are present",
			event: &storage.Event{
				Title:       "123",
				Description: "",
			},
		},
		{
			name: "all fields are there",
			event: &storage.Event{
				Title:       "test",
				Description: "test",
			},
		},
	}

	s := &Storage{
		ConnectionDB: newConnection(),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := s.CreateEvent(tc.event)
			if id == 0 {
				t.Error("id event == 0")
			}
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}
func TestGetByIdEvent(t *testing.T) {
	testCases := []struct {
		name string
		id   int
	}{
		{
			name: "get event by id == 1",
			id:   1,
		},
		{
			name: "get event by id == 2",
			id:   2,
		},
	}

	s := &Storage{
		ConnectionDB: newConnection(),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			event, err := s.GetByIdEvent(tc.id)
			fmt.Println(event)
			if event.ID != tc.id {
				t.Error("id event != 0")
			}
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestDeleteEventById(t *testing.T) {
	testCases := []struct {
		name string
		id   int
	}{
		{
			name: "1",
			id:   8,
		},
		{
			name: "2",
			id:   7,
		},
	}

	s := &Storage{
		ConnectionDB: newConnection(),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.DeleteEventById(tc.id)
			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func TestEditEvent(t *testing.T) {
	date := time.Now()
	testCases := []struct {
		name  string
		event *storage.Event
	}{
		{
			name: "all fields are there",
			event: &storage.Event{
				ID:            2,
				Title:         "test",
				Description:   "test",
				DateTimeStart: &date,
				DateTimeEnd:   &date,
			},
		},
		{
			name: "all fields are there",
			event: &storage.Event{
				ID:            4,
				Title:         "test2",
				Description:   "",
				DateTimeStart: &date,
				DateTimeEnd:   nil,
			},
		},
	}

	s := &Storage{
		ConnectionDB: newConnection(),
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.EditEvent(tc.event)

			if err != nil {
				t.Error(err.Error())
			}
		})
	}
}

func newConnection() *sql.DB {
	dsn := "postgres://postgres:qwer@localhost:5436/calendardb"
	storage, err := sql.Open("pgx", dsn)

	if err != nil {
		fmt.Println("Cannot open pgx driver: %w", err)
	}

	return storage
}
