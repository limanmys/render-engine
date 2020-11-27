package models

// ExtensionFileModel Structure of the extension data obj
type ExtensionFileModel struct {
	ID            string   `json:"id"`
	ExtensionID   string   `json:"extension_id"`
	Name          string   `json:"name"`
	Sha256sum     string   `json:"sha256sum"`
	ExtensionData []byte   `json:"extension_data"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	tableName     struct{} `pg:"extension_files"`
}
