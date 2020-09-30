package models

// License Structure of the personal extension license
type License struct {
	ID          string `json:"id"`
	Data        string `json:"data"`
	ExtensionID string `json:"extension_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
