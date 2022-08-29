package postgresql

import (
	"github.com/go-pg/pg/v10"
	"github.com/limanmys/render-engine/helpers"
	"github.com/limanmys/render-engine/models"
)

var db *pg.DB

// InitDB initialize database
func InitDB() {
	db = pg.Connect(&pg.Options{
		User:     helpers.DBUsername,
		Password: helpers.DBPassword,
		Database: helpers.DBName,
		Addr:     helpers.DBHost + ":" + helpers.DBPort,
	})

	//Thanks to library, we have to manually verify connection.
	data := struct {
		ID        string
		Migration string
		Batch     int
	}{}
	db.Query(&data, "select * from migrations limit 1")
	if data.ID == "" {
		panic("Postgresql sunusuna bağlanılamadı!")
	}
}

// GetUserData opens the database
func GetUserData(serverID string, extensionID string, userID string) (models.ServerModel, models.ExtensionModel, map[string]string) {
	return GetServer(serverID), GetExtension(extensionID), getSettings(userID, serverID, extensionID)
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
