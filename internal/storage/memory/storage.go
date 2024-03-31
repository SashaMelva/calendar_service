package memorystorage

import (
	"context"
	"database/sql"
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

func (s *Storage) CreateEvent(event *storage.Event) error {
	return nil
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
