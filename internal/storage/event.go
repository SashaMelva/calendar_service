package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Date         time.Duration
	Time         time.Duration
	Descriptioin string
}

// func SelectAll() {
// 	qwery := ""
// }
// func SelectById() {
// 	qwery := ""
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
