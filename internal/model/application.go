package model

type Application struct {
	ID       int64  `json:"id"`
	Company  string `json:"company"`
	Position string `json:"position"`
	Status   string `json:"status"`
}

type CreateApplicationRequest struct {
	Company  string `json:"company"`
	Position string `json:"position"`
	Status   string `json:"status"`
}
