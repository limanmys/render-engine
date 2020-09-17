package sqlite

import (
	"encoding/json"
	"renderer/src/helpers"
	"renderer/src/models"

	"github.com/mervick/aes-everywhere/go/aes256"
)

// GetWidget Get the id of the widget
func GetWidget(widgetID string) models.Widget {
	rows, _ := db.Query("SELECT * FROM widgets WHERE id=? LIMIT 1", widgetID)
	obj := models.Widget{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Title, &obj.UserID, &obj.Type, &obj.ExtensionID, &obj.ServerID, &obj.Function, &obj.Text, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order)
	rows.Close()
	return obj
}

//GetServer Retrieve server from id
func GetServer(serverID string) models.ServerModel {
	rows, _ := db.Query("SELECT * FROM servers WHERE id=? LIMIT 1", serverID)
	obj := models.ServerModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.ServerType, &obj.IPAddress, &obj.City, &obj.ControlPort, &obj.UserID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Os, &obj.Enabled, &obj.KeyPort)
	rows.Close()
	return obj
}

//GetExtension Retrieve extension from id
func GetExtension(extensionID string) models.ExtensionModel {
	rows, _ := db.Query("SELECT * FROM extensions WHERE id=? LIMIT 1", extensionID)
	obj := models.ExtensionModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Version, &obj.Icon, &obj.Service, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order, &obj.SslPorts, &obj.Issuer, &obj.Language, &obj.Support, &obj.Displays, &obj.Status)
	rows.Close()
	return obj
}

// GetExtensionFromName try to find extension id from it's name
func GetExtensionFromName(extensionName string) models.ExtensionModel {
	rows, _ := db.Query("SELECT * FROM extensions WHERE UPPER(NAME) LIKE UPPER(?) LIMIT 1", extensionName)
	obj := models.ExtensionModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Version, &obj.Icon, &obj.Service, &obj.CreatedAt, &obj.UpdatedAt, &obj.Order, &obj.SslPorts, &obj.Issuer, &obj.Language, &obj.Support, &obj.Displays, &obj.Status)
	rows.Close()
	return obj
}

// GetFuncPermissions Structure of the permissions
func GetFuncPermissions(userID string) []string {
	roleIDs := getRoleMapsFromUserID(userID)
	var permissions []string
	for _, roleID := range roleIDs {
		permissions = append(permissions, getFuncPermissionsFromMorphID(roleID)...)
	}

	permissions = append(permissions, getFuncPermissionsFromMorphID(userID)...)

	permissions = helpers.UniqueStrings(permissions)
	return permissions
}

// GetObjPermissions Structure of the permissions
func GetObjPermissions(userID string) []string {
	roleIDs := getRoleMapsFromUserID(userID)
	var permissions []string
	for _, roleID := range roleIDs {
		permissions = append(permissions, getObjPermissionsFromMorphID(roleID)...)
	}

	permissions = append(permissions, getObjPermissionsFromMorphID(userID)...)

	permissions = helpers.UniqueStrings(permissions)
	return permissions
}

func getRoleMapsFromUserID(userID string) []string {
	rows, _ := db.Query("SELECT * FROM role_users WHERE user_id=?", userID)
	var roleIDs []string
	for rows.Next() {
		obj := models.RoleUsers{}
		rows.Scan(&obj.ID, &obj.UserID, &obj.RoleID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Type)
		roleIDs = append(roleIDs, obj.RoleID)
	}
	rows.Close()
	return roleIDs
}

func getObjPermissionsFromMorphID(morphID string) []string {
	rows, _ := db.Query("SELECT * FROM permissions WHERE (morph_id=? and not type='function')", morphID)
	var permissions []string
	for rows.Next() {
		obj := models.Permission{}
		rows.Scan(&obj.ID, &obj.CreatedAt, &obj.UpdatedAt, &obj.Type, &obj.Key, &obj.Value, &obj.Extra, &obj.Blame, &obj.MorphID, &obj.MorphType)
		permissions = append(permissions, obj.Value)
	}
	rows.Close()
	return permissions
}

func getFuncPermissionsFromMorphID(morphID string) []string {
	rows, _ := db.Query("SELECT * FROM permissions WHERE (morph_id=? and type='function')", morphID)
	var permissions []string
	for rows.Next() {
		obj := models.Permission{}
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
		obj := models.SettingsModel{}
		rows.Scan(&obj.ID, &obj.ServerID, &obj.UserID, &obj.Name, &obj.Value, &obj.CreatedAt, &obj.UpdatedAt)
		obj.Value = aes256.Decrypt(obj.Value, decryptionKey)
		results[obj.Name] = obj.Value
	}
	rows.Close()
	return results
}

func getToken(token string) (models.TokenModel, error) {
	rows, err := db.Query("SELECT * FROM tokens WHERE token=? LIMIT 1", token)
	if err != nil {
		return models.TokenModel{}, err
	}
	obj := models.TokenModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.UserID, &obj.Token, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

func getAccessToken(token string) (models.AccessToken, error) {
	rows, err := db.Query("SELECT * FROM access_tokens WHERE token=? LIMIT 1", token)
	if err != nil {
		return models.AccessToken{}, err
	}
	obj := models.AccessToken{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.UserID, &obj.LastUsedAt, &obj.LastUsedIP, &obj.Token, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

// GetLicense Structure of the license object
func GetLicense(extensionID string) (models.License, error) {
	rows, err := db.Query("SELECT * FROM licenses WHERE extension_id=? LIMIT 1", extensionID)
	if err != nil {
		return models.License{}, err
	}
	obj := models.License{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Data, &obj.ExtensionID, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	return obj, nil
}

// GetUser Retrieve user data from id
func GetUser(userID string) models.UserModel {
	rows, _ := db.Query("SELECT * FROM users WHERE id=? LIMIT 1", userID)
	obj := models.UserModel{}
	rows.Next()
	rows.Scan(&obj.ID, &obj.Name, &obj.Email, &obj.Password, &obj.Status, &obj.LastLoginAt, &obj.RememberToken, &obj.LastLoginIP, &obj.CreatedAt, &obj.UpdatedAt, &obj.ForceChange, &obj.ObjectGUID, &obj.AuthType)
	rows.Close()
	obj.Password = ""
	obj.RememberToken = ""
	obj.ObjectGUID = ""
	return obj
}

//GetServerKey Retrieve the user key.
func GetServerKey(userID string, serverID string) (string, string, string, models.ServerKey) {
	decryptionKey := helpers.AppKey + userID + serverID
	rows, _ := db.Query("SELECT * FROM server_keys WHERE (user_id=? AND server_id=? )", userID, serverID)
	rows.Next()
	obj := models.ServerKey{}
	rows.Scan(&obj.ID, &obj.Type, &obj.Data, &obj.ServerID, &obj.UserID, &obj.CreatedAt, &obj.UpdatedAt)
	rows.Close()
	if obj.Data == "" {
		return "", "", "", models.ServerKey{}
	}
	type keyData struct {
		ClientUsername string `json:"clientUsername"`
		ClientPassword string `json:"clientPassword"`
		KeyPort        string `json:"key_port"`
	}
	var key keyData

	json.Unmarshal([]byte(obj.Data), &key)
	return aes256.Decrypt(key.ClientUsername, decryptionKey), aes256.Decrypt(key.ClientPassword, decryptionKey), key.KeyPort, obj
}
