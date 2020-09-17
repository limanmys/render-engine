package models

// Permission Structure of the permissions
type Permission struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Type      string `json:"type"`
	Key       string `json:"key"`
	Value     string `json:"value"`
	Extra     string `json:"extra"`
	Blame     string `json:"blame"`
	MorphID   string `json:"morph_id"`
	MorphType string `json:"morph_type"`
}
