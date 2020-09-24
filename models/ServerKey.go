package models

// ServerKey Structure of the server keys
type ServerKey struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Data      string `json:"data"`
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
