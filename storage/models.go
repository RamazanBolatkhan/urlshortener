package storage

import "time"

type URL struct {
	OriginalUrl string    `json:"original_url"`
	Alias       string    `json:"alias"`
	CreatedTime time.Time `json:"created_time"`
}
