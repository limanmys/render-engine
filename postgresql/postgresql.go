package postgresql

import (
	"io/ioutil"

	"github.com/go-pg/pg/v10"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/models"
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
	return GetServer(serverID), GetExtension(extensionID), getSettings(userID, serverID)
}

//ReadSettings Read Settings from Configuration File
func ReadSettings() {
	var settings []models.SystemSettingsModel
	_ = db.Model(&settings).ForEach(func(setting models.SystemSettingsModel) error {
		switch setting.Key {
		case "APP_KEY":
			helpers.AppKey = setting.Data
			break
		case "GO_KEY":
			helpers.AuthKey = setting.Data
			break
		case "LIMAN_IP":
			helpers.LimanIP = setting.Data
			break
		case "LIMAN_RESTRICTED":
			helpers.LimanIP = setting.Data
			break
		case "SSL_PRIVATE_KEY":
			d1 := []byte(setting.Data)
			ioutil.WriteFile(helpers.CertsPath+"liman.key", d1, 0600)
			break
		case "SSL_PUBLIC_KEY":
			d1 := []byte(setting.Data)
			ioutil.WriteFile(helpers.CertsPath+"liman.crt", d1, 0600)
			break
		}

		return nil
	})
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
