package models

// SystemSettingsModel Structure of the user system settings obj
type SystemSettingsModel struct {
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	Data      string   `json:"data"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	tableName struct{} `pg:"system_settings"`
}
