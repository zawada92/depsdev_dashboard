package model

import "time"

type Dependency struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	OpenSSFScore float64   `json:"openssf_score"`
	LastUpdated  time.Time `json:"last_updated"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
