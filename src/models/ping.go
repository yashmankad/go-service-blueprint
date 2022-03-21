package models

// PingResponse is the server response for the ping API endpoint
type PingResponse struct {
	Message string `json:"message"`
}
