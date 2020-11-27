package replications

import (
	"encoding/base64"
	"errors"
	"io/ioutil"

	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/postgresql"
)

func DownloadExtension(extensionID string) (string, error) {
	obj := postgresql.GetExtensionFile(extensionID)

	if obj.Sha256sum == "" {
		return "", errors.New("Eklenti bulunamadı!")
	}
	sDec, _ := base64.StdEncoding.DecodeString(string(obj.ExtensionData))

	path := "/tmp/" + obj.Name

	err := ioutil.WriteFile(path, sDec, 0644)

	if err != nil {
		return "", err
	}

	sha256sum, _ := helpers.ExecuteCommand("sha256sum " + path + " | awk '{ print $1 }'")

	if sha256sum != obj.Sha256sum {
		return "", errors.New("İndirilen dosyanın hashi doğrulanamadı! " + obj.Sha256sum)
	}

	return path, nil
}
