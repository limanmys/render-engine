package postgresql

import (
	"github.com/go-pg/pg/v10"
	"github.com/limanmys/go/models"
)

var db *pg.DB


// GetUserData opens the database
func GetUserData(serverID string, extensionID string, userID string) (models.ServerModel, models.ExtensionModel, map[string]string) {
	return GetServer(serverID), GetExtension(extensionID), getSettings(userID, serverID)
}

// InitDB initialize database
func InitDB() {
	db = pg.Connect(&pg.Options{
		User: "liman",
		Password: "XczvmFNBIECYiEnVD9s1NURhs",
		Database: "liman",
	})
}