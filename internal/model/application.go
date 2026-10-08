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

type UpdateApplicationRequest struct {
	Company   OptionalString `json:"company"`
	Position  OptionalString `json:"position"`
	Status    OptionalString `json:"status"`
	AppliedAt OptionalString `json:"applied_at"`
	Notes     OptionalString `json:"notes"`
}

type ListApplicationsRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Status   string `form:"status"`
	Company  string `form:"company"`
}

type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type ListApplicationsResponse struct {
	Data       []Application `json:"data"`
	Pagination Pagination    `json:"pagination"`
}
