package models

// SettingsModel Structure of the user settings obj
type SettingsModel struct {
	ID        string `json:"id"`
	ServerID  string `json:"server_id"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
