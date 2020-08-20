package sqlite

import (
	"renderer/src/helpers"

	"github.com/mervick/aes-everywhere/go/aes256"
)

// GetWidget Get the id of the widget
func GetWidget(widgetID string) Widget {
	rows, _ := db.Query("SELECT * FROM widgets WHERE id=? LIMIT 1", widgetID)
	obj := Widget{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Title, &obj.UserID, &obj.Type, &obj.ExtensionID, &obj.ServerID, &obj.Function, &obj.Text, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order)
	rows.Close()
	return obj
}

func getServer(serverID string) ServerModel {
	rows, _ := db.Query("SELECT * FROM servers WHERE id=? LIMIT 1", serverID)
	obj := ServerModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.ServerType, &obj.IPAddress, &obj.City, &obj.ControlPort, &obj.UserID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Os, &obj.Enabled, &obj.KeyPort)
	rows.Close()
	return obj
}

func getExtension(extensionID string) ExtensionModel {
	rows, _ := db.Query("SELECT * FROM extensions WHERE id=? LIMIT 1", extensionID)
	obj := ExtensionModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Version, &obj.Icon, &obj.Service, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order, &obj.SslPorts, &obj.Issuer, &obj.Language, &obj.Support, &obj.Displays, &obj.Status)
	rows.Close()
	return obj
}

// GetPermissions Structure of the permissions
func GetPermissions(userID string) []string {
	roleIDs := getRoleMapsFromUserID(userID)
	var permissions []string
	for _, roleID := range roleIDs {
		permissions = append(permissions, getPermissionsFromMorphID(roleID)...)
	}

	permissions = append(permissions, getPermissionsFromMorphID(userID)...)

	permissions = helpers.UniqueStrings(permissions)
	return permissions
}

func getRoleMapsFromUserID(userID string) []string {
	rows, _ := db.Query("SELECT * FROM role_users WHERE user_id=?", userID)
	var roleIDs []string
	for rows.Next() {
		obj := RoleUsers{}
		rows.Scan(&obj.ID, &obj.UserID, &obj.RoleID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Type)
		roleIDs = append(roleIDs, obj.RoleID)
	}
	rows.Close()
	return roleIDs
}

func getPermissionsFromMorphID(morphID string) []string {
	rows, _ := db.Query("SELECT * FROM permissions WHERE (morph_id=? and type='function')", morphID)
	var permissions []string
	for rows.Next() {
		obj := Permission{}
		rows.Scan(&obj.ID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Type, &obj.Key, &obj.Value, &obj.Extra, &obj.Blame, &obj.MorphID, &obj.MorphType)
		permissions = append(permissions, obj.Extra)
	}
	rows.Close()
	return permissions
}

func getSettings(userID string, serverID string) map[string]string {
	rows, _ := db.Query("SELECT * FROM user_settings WHERE (user_id=? AND server_id=? )", userID, serverID)
	results := make(map[string]string)
	decryptionKey := helpers.AppKey + userID + serverID
	for rows.Next() {
		obj := SettingsModel{}
		rows.Scan(&obj.ID, &obj.ServerID, &obj.UserID, &obj.Name, &obj.Value, &obj.CreatedAt, &obj.UpdatedAt)
		obj.Value = aes256.Decrypt(obj.Value, decryptionKey)
		results[obj.Name] = obj.Value
	}
	rows.Close()
	return results
}

func getToken(token string) (TokenModel, error) {
	rows, err := db.Query("SELECT * FROM tokens WHERE token=? LIMIT 1", token)
	if err != nil {
		return TokenModel{}, err
	}
	obj := TokenModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.UserID, &obj.Token, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

func getAccessToken(token string) (AccessToken, error) {
	rows, err := db.Query("SELECT * FROM access_tokens WHERE token=? LIMIT 1", token)
	if err != nil {
		return AccessToken{}, err
	}
	obj := AccessToken{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.UserID, &obj.LastUsedAt, &obj.LastUsedIP, &obj.Token, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

// GetLicense Structure of the license object
func GetLicense(extensionID string) (License, error) {
	rows, err := db.Query("SELECT * FROM licenses WHERE extension_id=? LIMIT 1", extensionID)
	if err != nil {
		return License{}, err
	}
	obj := License{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Data, &obj.ExtensionID, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

// GetUser Retrieve user data from id
func GetUser(userID string) UserModel {
	rows, _ := db.Query("SELECT * FROM users WHERE id=? LIMIT 1", userID)
	obj := UserModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Email, &obj.Password, &obj.Status, &obj.LastLoginAt, &obj.RememberToken, &obj.LastLoginIP, &obj.CreatedAt, &obj.UpdatedAt, &obj.ForceChange, &obj.ObjectGUID, &obj.AuthType)
	rows.Close()
	obj.Password = ""
	obj.RememberToken = ""
	obj.ObjectGUID = ""
	return obj
}
