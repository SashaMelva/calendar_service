package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Date         time.Duration
	Time         time.Time
	Descriptioin string
}
