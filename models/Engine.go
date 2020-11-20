package models

// EngineModel Structure of the engine obj
type EngineModel struct {
	ID        string   `json:"id"`
	Token     string   `json:"token" pg:"token"`
	MachineID string   `json:"machine_id" pg:"machine_id"`
	IPAddress string   `json:"ip_address" pg:"ip_address"`
	Port      int      `json:"port"`
	Enabled   bool     `json:"enabled"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	tableName struct{} `pg:"go_engines"`
}
