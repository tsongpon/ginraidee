package model

import "time"

type SearchHistory struct {
	ID      string
	UserID  string
	Keyword string
	Time    time.Time
}
