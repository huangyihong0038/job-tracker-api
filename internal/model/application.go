package model

import "time"

type Application struct {
	ID        int64     `json:"id"`
	Company   string    `json:"company"`
	Position  string    `json:"position"`
	Status    string    `json:"status"`
	AppliedAt *string   `json:"applied_at"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateApplicationRequest struct {
	Company   string  `json:"company"`
	Position  string  `json:"position"`
	Status    string  `json:"status"`
	AppliedAt *string `json:"applied_at"`
	Notes     string  `json:"notes"`
}
