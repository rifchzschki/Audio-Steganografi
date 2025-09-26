package models

type MessageResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}