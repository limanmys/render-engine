package postgresql

import (
	"fmt"
	"github.com/limanmys/go/models"
)

func GetWidget(widgetID string) models.Widget {
	widget := &models.Widget{ID: widgetID}
	_ = db.Model(widget).WherePK().First()
	return *widget
}

//GetServer Retrieve server from id
func GetServer(serverID string) models.ServerModel {
	server := &models.ServerModel{ID: serverID}
	_ = db.Model(server).WherePK().Select()
	return *server
}

//GetExtension Retrieve extension from id
func GetExtension(extensionID string) models.ExtensionModel {
	extension := &models.ExtensionModel{ID: extensionID}
	_ = db.Model(extension).WherePK().First()
	return *extension
}

// GetExtensionFromName try to find extension id from it's name
func GetExtensionFromName(extensionName string) models.ExtensionModel {
	extension := &models.ExtensionModel{}
	_ = db.Model(extension).Where("name = ?", extensionName).First()
	return *extension
}

//// GetFuncPermissions Structure of the permissions
//func GetFuncPermissions(userID string) []string {
//	roleIDs := getRoleMapsFromUserID(userID)
//	var permissions []string
//	for _, roleID := range roleIDs {
//		permissions = append(permissions, getFuncPermissionsFromMorphID(roleID)...)
//	}
//
//	permissions = append(permissions, getFuncPermissionsFromMorphID(userID)...)
//
//	permissions = helpers.UniqueStrings(permissions)
//	return permissions
//}

//// GetObjPermissions Structure of the permissions
//func GetObjPermissions(userID string) []string {
//	roleIDs := getRoleMapsFromUserID(userID)
//	var permissions []string
//	for _, roleID := range roleIDs {
//		permissions = append(permissions, getObjPermissionsFromMorphID(roleID)...)
//	}
//
//	permissions = append(permissions, getObjPermissionsFromMorphID(userID)...)
//
//	permissions = helpers.UniqueStrings(permissions)
//	return permissions
//}

func getRoleMapsFromUserID(userID string) []string {
	var roles []models.RoleUsers
	var roleIDs []string
	_ = db.Model(&roles).Where("user_id = ?", userID).ForEach(func(role models.RoleUsers) error {
		roleIDs = append(roleIDs, role.RoleID)
		return nil
	})
	return roleIDs
}


func getSettings(userID string, serverID string) map[string]string {
	var settings []models.SettingsModel
	_, _ = db.Model(&settings).Query("SELECT * FROM user_settings WHERE (user_id=? AND server_id=? )", userID, serverID)
	fmt.Println(settings)
	return map[string]string{}
	/*results := make(map[string]string)
	decryptionKey := helpers.AppKey + userID + serverID
	for key, data := range settings {

	}
	for rows.Next() {
		obj := models.SettingsModel{}
		rows.Scan(&obj.ID, &obj.ServerID, &obj.UserID, &obj.Name, &obj.Value, &obj.CreatedAt, &obj.UpdatedAt)
		obj.Value = aes256.Decrypt(obj.Value, decryptionKey)
		results[obj.Name] = obj.Value
	}
	rows.Close()
	return results*/
}