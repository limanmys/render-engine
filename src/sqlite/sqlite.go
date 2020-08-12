package sqlite

import (
	"database/sql"
	"renderer/src/helpers"

	//Sqlite3
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// GetUserData opens the database
func GetUserData(serverID string, extensionID string, userID string) (ServerModel, ExtensionModel, []SettingsModel) {
	return getServer(serverID), getExtension(extensionID), getSettings(userID, serverID)
}

// GetUserIDFromToken Find token from token
func GetUserIDFromToken(tokenID string) string {
	token, err := getToken(tokenID)
	if err != nil {
		return ""
	}
	return token.UserID
}

// InitDB inialize database
func InitDB() {
	temp, err := sql.Open("sqlite3", "/liman/database/liman.sqlite")
	if err != nil {
		helpers.Abort(err.Error())
	}
	db = temp
}
