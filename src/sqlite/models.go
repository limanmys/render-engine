package sqlite

import (
	"renderer/src/helpers"

	"github.com/mervick/aes-everywhere/go/aes256"
)

// ServerModel Structure of the server obj
type ServerModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ServerType  string `json:"type"`
	IPAddress   string `json:"ip_address"`
	City        string `json:"city"`
	ControlPort string `json:"control_port"`
	UserID      string `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Os          string `json:"os"`
	Enabled     string `json:"enabled"`
	KeyPort     int    `json:"key_port"`
}

// ExtensionModel Structure of the extension obj
type ExtensionModel struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Version   string `json:"version"`
	Icon      string `json:"icon"`
	Service   string `json:"service"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Order     int    `json:"order"`
	SslPorts  string `json:"sslPorts"`
	Issuer    string `json:"issuer"`
	Language  string `json:"language"`
	Support   string `json:"support"`
	Displays  string `json:"displays"`
	Status    string `json:"status"`
}

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

// TokenModel Structure of the tokens table
type TokenModel struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Token     string `json:"token"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// UserModel Structure of the users table
type UserModel struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Status        int    `json:"status"`
	LastLoginAt   string `json:"last_login_at"`
	RememberToken string `json:"remember_token"`
	LastLoginIP   string `json:"last_login_ip"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	ForceChange   int    `json:"forceChange"`
	ObjectGUID    string `json:"objectguid"`
	AuthType      string `json:"auth_type"`
}

func getServer(serverID string) ServerModel {
	rows, _ := db.Query("SELECT * FROM servers WHERE id='" + serverID + "' LIMIT 1")
	obj := ServerModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.ServerType, &obj.IPAddress, &obj.City, &obj.ControlPort, &obj.UserID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Os, &obj.Enabled, &obj.KeyPort)
	return obj
}

func getExtension(extensionID string) ExtensionModel {
	rows, _ := db.Query("SELECT * FROM extensions WHERE id='" + extensionID + "' LIMIT 1")
	obj := ExtensionModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Version, &obj.Icon, &obj.Service, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order, &obj.SslPorts, &obj.Issuer, &obj.Language, &obj.Support, &obj.Displays, &obj.Status)
	return obj
}

func getSettings(userID string, serverID string) []SettingsModel {
	rows, _ := db.Query("SELECT * FROM user_settings WHERE (user_id='" + userID + "' AND server_id='" + serverID + "' )")
	results := []SettingsModel{}
	decryptionKey := helpers.AppKey + userID + serverID
	for rows.Next() {
		obj := SettingsModel{}
		rows.Scan(&obj.ID, &obj.ServerID, &obj.UserID, &obj.Name, &obj.Value, &obj.CreatedAt, &obj.UpdatedAt)
		obj.Value = aes256.Decrypt(obj.Value, decryptionKey)
		results = append(results, obj)
	}
	return results
}

func getToken(token string) (TokenModel, error) {
	rows, err := db.Query("SELECT * FROM tokens WHERE token='" + token + "' LIMIT 1")
	if err != nil {
		return TokenModel{}, err
	}
	obj := TokenModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.UserID, &obj.Token, &obj.CreatedAt, &obj.UpdatedAt)
	return obj, nil
}

// GetUser Retrieve user data from id
func GetUser(userID string) UserModel {
	rows, _ := db.Query("SELECT * FROM users WHERE id='" + userID + "' LIMIT 1")
	obj := UserModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Email, &obj.Password, &obj.Status, &obj.LastLoginAt, &obj.RememberToken, &obj.LastLoginIP, &obj.CreatedAt, &obj.UpdatedAt, &obj.ForceChange, &obj.ObjectGUID, &obj.AuthType)
	obj.Password = ""
	obj.RememberToken = ""
	obj.ObjectGUID = ""
	return obj
}
