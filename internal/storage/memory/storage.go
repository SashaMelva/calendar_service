package memorystorage

import (
	"context"
	"database/sql"
	"sync"
	"time"

	storage "github.com/SashaMelva/calendar_service/internal/storage"
)

type Storage struct {
	Ctx          context.Context
	ConnectionDB *sql.DB
	sync.RWMutex //nolint:unused
}

func New(connection *sql.DB) *Storage {
	return &Storage{
		ConnectionDB: connection,
	}
}

func (s *Storage) GetAllEvents() ([]storage.Event, error) {
	events := []storage.Event{}
	query := `select * from events`
	rows, err := s.ConnectionDB.QueryContext(s.Ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := storage.Event{}

		if err := rows.Scan(event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) GetByIdEvent(id int) (*storage.Event, error) {
	event := &storage.Event{}
	query := `select id, title, date_time_start, date_time_end, description from events where id = $1`
	row := s.ConnectionDB.QueryRow(query, id)

	err := row.Scan(
		&event.ID,
		&event.Title,
		&event.DateTimeStart,
		&event.DateTimeEnd,
		&event.Description,
	)

	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return event, err
	}

	return event, nil
}

func (s *Storage) CreateEvent(event *storage.Event) (int, error) {
	var eventId int
	query := `insert into events(title, description, date_time_start, date_time_end) values($1, $2, $3, $4) RETURNING id`
	result := s.ConnectionDB.QueryRow(query, event.Title, event.Description, event.DateTimeStart, event.DateTimeEnd) // sql.Result
	err := result.Scan(&eventId)

	if err != nil {
		return 0, err
	}

	return eventId, nil
}

func (s *Storage) DeleteEventById(id int) error {
	query := `delete from events where id = $1`
	_, err := s.ConnectionDB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) EditEvent(event *storage.Event) error {
	query := `update events set title=$1, description=$2, date_time_start=$3, date_time_end=$4 where id=$5`
	_, err := s.ConnectionDB.Exec(query, event.Title, event.Description, event.DateTimeStart, event.DateTimeEnd, event.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListEventsDateForPeriod(dateStart, dateEnd *time.Time) ([]storage.Event, error) {
	events := []storage.Event{}
	query := `select * from test  where date >= date_time_start $1 and date < date_time_start $2`
	rows, err := s.ConnectionDB.QueryContext(s.Ctx, query, dateStart, dateEnd)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		event := storage.Event{}

		if err := rows.Scan(event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
