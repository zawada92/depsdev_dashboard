package model

import "time"

type Dependency struct {
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	OpenSSFScore float64   `json:"openssf_score"`
	LastUpdated  time.Time `json:"last_updated"`
}
type PatchDependencyRequest struct {
	Version      *string  `json:"version,omitempty"`
	OpenSSFScore *float64 `json:"openssf_score,omitempty"`
}
