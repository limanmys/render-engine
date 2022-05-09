package postgresql

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/models"
	"github.com/mervick/aes-everywhere/go/aes256"
)

func GetWidget(widgetID string) models.Widget {
	widget := &models.Widget{ID: widgetID}
	_ = db.Model(widget).WherePK().First()
	return *widget
}

//GetServer Retrieve server from id
func GetServer(serverID string) models.ServerModel {
	server := &models.ServerModel{ID: serverID}
	db.Model(server).WherePK().Select()
	return *server
}

//GetExtension Retrieve extension from id
func GetExtension(extensionID string) models.ExtensionModel {
	extension := &models.ExtensionModel{ID: extensionID}
	db.Model(extension).WherePK().First()
	return *extension
}

// GetExtensionFromName try to find extension id from it's name
func GetExtensionFromName(extensionName string) models.ExtensionModel {
	extension := &models.ExtensionModel{}
	db.Model(extension).Where("name = ?", extensionName).First()
	return *extension
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

// GetVariables Structure of the permissions
func GetVariables(userID string) map[string]string {
	roleIDs := getRoleMapsFromUserID(userID)
	permissions := make(map[string]string)
	for _, roleID := range roleIDs {
		permissions = helpers.MergeStringMaps(permissions, getVariablesFromMorphID(roleID))
	}
	permissions = helpers.MergeStringMaps(permissions, getVariablesFromMorphID(userID))

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
	var roles []models.RoleUsers
	var roleIDs []string
	_ = db.Model(&roles).Where("user_id = ?", userID).ForEach(func(role models.RoleUsers) error {
		roleIDs = append(roleIDs, role.RoleID)
		return nil
	})
	return roleIDs
}

func getObjPermissionsFromMorphID(morphID string) []string {
	var permission []models.Permission
	var permissions []string
	_ = db.Model(&permission).Where("morph_id=? and not type='function'", morphID).ForEach(func(permission models.Permission) error {
		permissions = append(permissions, permission.Value)
		return nil
	})
	return permissions
}

func getFuncPermissionsFromMorphID(morphID string) []string {
	var permission []models.Permission
	var permissions []string
	_ = db.Model(&permission).Where("morph_id=? and type='function'", morphID).ForEach(func(permission models.Permission) error {
		permissions = append(permissions, permission.Extra)
		return nil
	})
	return permissions
}

func getVariablesFromMorphID(morphID string) map[string]string {
	var permission []models.Permission
	permissions := make(map[string]string)
	_ = db.Model(&permission).Where("morph_id=? and type='variable'", morphID).ForEach(func(permission models.Permission) error {
		permissions[permission.Key] = permission.Value
		return nil
	})
	return permissions
}

func getSettings(userID string, serverID string, extensionID string) map[string]string {
	extensionJSON, err := GetExtensionJSON(extensionID)

	if err != nil {
		return make(map[string]string)
	}

	var globalVars []string
	for _, setting := range extensionJSON["database"].([]interface{}) {
		isGlobal := setting.(map[string]interface{})["global"]
		if isGlobal != nil && isGlobal.(bool) {
			globalVars = append(globalVars, setting.(map[string]interface{})["variable"].(string))
		}
	}

	var settings []models.SettingsModel
	results := make(map[string]string)

	decryptionKey := helpers.AppKey + userID + serverID
	_ = db.Model(&settings).Where("user_id=? AND server_id=?", userID, serverID).ForEach(func(setting models.SettingsModel) error {
		setting.Value = aes256.Decrypt(setting.Value, decryptionKey)
		results[setting.Name] = setting.Value
		return nil
	})

	_ = db.Model(&settings).Where("name IN (?) AND server_id=?", pg.In(globalVars), serverID).ForEach(func(setting models.SettingsModel) error {
		if _, ok := results[setting.Name]; ok {
			return nil
		}
		setting.Value = aes256.Decrypt(setting.Value, helpers.AppKey+setting.UserID+serverID)
		results[setting.Name] = setting.Value
		return nil
	})
	return results
}

func getToken(token string) (models.TokenModel, error) {
	object := &models.TokenModel{}
	err := db.Model(object).Where("token = ?", token).First()
	return *object, err
}

func getAccessToken(token string) (models.AccessToken, error) {
	object := &models.AccessToken{}
	err := db.Model(object).Where("token = ?", token).First()
	return *object, err
}

// GetLicense Structure of the license object
func GetLicense(extensionID string) (models.License, error) {
	object := &models.License{}
	err := db.Model(object).Where("extension_id = ?", extensionID).First()
	return *object, err
}

// GetUser Retrieve user data from id
func GetUser(userID string) models.UserModel {
	user := &models.UserModel{ID: userID}
	err := db.Model(user).WherePK().First()
	if err != nil {
		panic(err)
	}
	user.Password = ""
	user.RememberToken = ""
	user.ObjectGUID = ""
	return *user
}

//GetServerKey Retrieve the user key.
func GetServerKey(userID string, serverID string) (string, string, string, models.ServerKey) {
	object := &models.ServerKey{}
	decryptionKey := helpers.AppKey + userID + serverID

	db.Model(object).Where("user_id=? AND server_id=?", userID, serverID).First()

	if object.Data == "" {
		server := &models.ServerModel{}
		db.Model(server).Where("id=?", serverID).First()

		if server.SharedKey == 1 {
			db.Model(object).Where("server_id=?", serverID).First()
			decryptionKey = helpers.AppKey + object.UserID + serverID
		} else {
			return "", "", "", models.ServerKey{}
		}
	}

	type keyData struct {
		ClientUsername string `json:"clientUsername"`
		ClientPassword string `json:"clientPassword"`
		KeyPort        string `json:"key_port"`
	}
	var key keyData

	json.Unmarshal([]byte(object.Data), &key)
	return aes256.Decrypt(key.ClientUsername, decryptionKey), aes256.Decrypt(key.ClientPassword, decryptionKey), key.KeyPort, *object
}

//GetExtensionJSON Get extension json file with extension id
func GetExtensionJSON(extensionID string) (map[string]interface{}, error) {
	extension := GetExtension(extensionID)
	extensionJSONFile, err := ioutil.ReadFile("/liman/extensions/" + strings.ToLower(extension.Name) + "/db.json")

	if err != nil {
		return nil, err
	}

	var jsonMap map[string]interface{}
	_ = json.Unmarshal(extensionJSONFile, &jsonMap)

	return jsonMap, nil
}
