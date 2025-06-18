package domain

import "time"

type Users struct {
	Id           string     `json:"id"`
	Username     string     `json:"username"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	LastAccessed *time.Time `json:"last_accessed"`
}
