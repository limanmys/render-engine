package sqlite

import (
	"database/sql"
	"renderer/src/helpers"

	//Sqlite3
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

// GetUserData opens the database
func GetUserData(serverID string, extensionID string, userID string) (ServerModel, ExtensionModel, map[string]string) {
	return GetServer(serverID), GetExtension(extensionID), getSettings(userID, serverID)
}

// GetUserIDFromToken Find token from token
func GetUserIDFromToken(tokenID string) string {
	token, err := getToken(tokenID)
	if err != nil {
		return ""
	}
	return token.UserID
}

// GetUserIDFromLimanToken Find token from liman token (personal access key)
func GetUserIDFromLimanToken(tokenID string) string {
	token, err := getAccessToken(tokenID)
	if err != nil {
		return ""
	}
	return token.UserID
}

// InitDB inialize database
func InitDB() {
	temp, err := sql.Open("sqlite3", "/liman/database/liman.sqlite?cache=shared&mode=rwc")
	if err != nil {
		helpers.Abort(err.Error())
	}
	db = temp
}
