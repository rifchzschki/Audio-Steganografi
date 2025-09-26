package models

// MessageResponse merepresentasikan struktur data yang dikirim sebagai JSON response.
type MessageResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}