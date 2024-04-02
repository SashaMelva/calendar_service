package memorystorage

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

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
	query := `select * from events where id = $1`
	row := s.ConnectionDB.QueryRowContext(s.Ctx, query)

	err := row.Scan(event)

	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return event, err
	}

	return event, nil
}

func (s *Storage) CreateEvent(event *storage.Event) (int64, error) {
	fmt.Println("sql")
	query := `insert into events(title, descriptioin) values($1, $2)`
	result, err := s.ConnectionDB.Exec(query, event.Title, event.Descriptioin) // sql.Result

	if err != nil {
		return 0, err
	}
	fmt.Println("sql2")
	eventId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return eventId, nil
}
func (s *Storage) DeleteEvent(id int) error {
	return nil
}
func (s *Storage) EditEvent(event *storage.Event) error {
	return nil
}

// func SelectAll() {
// 	qwery := "select * from events"
// }

// func SelectById() {
// 	query := `
// 			select *
// 			from events
// 			where id = $1
// 			`
// 	rows, err := db.QueryContext(ctx, query, id)

// 	if err == sql.ErrNoRows {
// 		// строки не найдено
// 	} else if err != nil {
// 		// "настоящая" ошибка
// 	}

// 	defer rows.Close()

// 	for rows.Next() {
// 		var id int64
// 		var title, descr string
// 		if err := rows.Scan(&id, &title, &descr); err != nil {
// 			// ошибка сканирования
// 		}
// 		// обрабатываем строку
// 		fmt.Printf("%d %s %s\n", id, title, descr)
// 	}
// 	if err := rows.Err(); err != nil {
// 		// ошибка при получении результатов
// 	}
// }

// func Create() {
// 	qwery := ""
// }

// func Delete() {
// 	qwery := ""
// }

// func Edit() {
// 	qwery := ""
// }
