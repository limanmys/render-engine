package postgresql

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/models"
	"github.com/mervick/aes-everywhere/go/aes256"
)

func SetExtensionDb(value string, target string, serverID string, isGlobal bool, userID string) string {
	var found bool = false
	var settings []models.SettingsModel
	if isGlobal {
		_ = db.Model(&settings).Where("name = ? AND server_id = ?", target, serverID).ForEach(func(setting models.SettingsModel) error {
			key := helpers.AppKey + setting.UserID + setting.ServerID
			setting.Value = aes256.Encrypt(value, key)
			setting.UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			_, _ = db.Model(&setting).WherePK().Update()
			return nil
		})
	} else {
		_ = db.Model(&settings).Where("name = ? AND user_id = ? AND server_id = ?", target, userID, serverID).ForEach(func(setting models.SettingsModel) error {
			key := helpers.AppKey + setting.UserID + setting.ServerID
			setting.Value = aes256.Encrypt(value, key)
			setting.UpdatedAt = time.Now().Format(time.RFC3339)
			found = true
			_, _ = db.Model(&setting).WherePK().Update()
			return nil
		})
	}

	if !found {
		log.Println("insert")
		uuid, _ := uuid.NewUUID()
		setting := &models.SettingsModel{
			ID:        uuid.String(),
			ServerID:  serverID,
			UserID:    userID,
			Name:      target,
			Value:     aes256.Encrypt(value, helpers.AppKey+userID+serverID),
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		_, _ = db.Model(setting).Insert()
	}

	return value
}
