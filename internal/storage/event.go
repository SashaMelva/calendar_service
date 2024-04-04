package storage

import "time"

type MyTime *time.Time

type Event struct {
	ID            int        `json:"id" db:"id"`
	Title         string     `json:"title" db:"title"`
	DateTimeStart *time.Time `json:"date_time_start" db:"date_time_start"`
	DateTimeEnd   *time.Time `json:"date_time_end" db:"date_time_end"`
	Description   string     `json:"description" db:"description"`
}

func Time(myTime *time.Time) *string {
	if myTime != nil {
		time := myTime.Format("15:04:05")
		return &time
	}

	return nil
}

func Date(myTime *time.Time) *string {
	if myTime != nil {
		date := myTime.Format("2006.01.02")
		return &date
	}

	return nil
}
