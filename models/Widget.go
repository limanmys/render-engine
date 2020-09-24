package models

// Widget Structure of the widget
type Widget struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Title       string `json:"title"`
	UserID      string `json:"user_id"`
	Type        string `json:"type"`
	ExtensionID string `json:"extension_id"`
	ServerID    string `json:"serrver_id"`
	Function    string `json:"function"`
	Text        string `json:"text"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Order       int    `json:"order"`
}
