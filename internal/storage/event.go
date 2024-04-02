package storage

import "time"

type Event struct {
	ID           string        `json:"id" db:"id"`
	Title        string        `json:"title" db:"title"`
	Date         time.Duration `json:"date" db:"date"`
	Time         time.Time     `json:"time" db:"time"`
	Descriptioin string        `json:"descriptioin" db:"id"`
}
