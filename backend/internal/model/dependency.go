package model

import "time"

type Dependency struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	OpenSSFScore float64   `json:"openssf_score"`
	LastUpdated  time.Time `json:"last_updated"`
}
