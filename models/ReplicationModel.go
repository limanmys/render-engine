package models

// ReplicationModel Structure of the server obj
type ReplicationModel struct {
	ID        string   `json:"id"`
	Key       string   `json:"key"`
	MachineID string   `json:"machine_id"`
	Log       string   `json:"log"`
	Completed bool     `json:"completed"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	tableName struct{} `pg:"replications"`
}
