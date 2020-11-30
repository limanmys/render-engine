package replications

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/postgresql"
)

//DNS Dns Replication Loop
func DNS() {
	log.Println("Checking DNS updates from Liman")
	setting := postgresql.GetSystemSetting("SYSTEM_DNS")
	if setting.ID == "" {
		return
	}

	replication := postgresql.GetReplication("SYSTEM_DNS")

	if helpers.IsNewer(replication.UpdatedAt, setting.UpdatedAt) {
		log.Println("Dns already up to date.")
		return
	}

	var dns = struct {
		Server1 string
		Server2 string
		Server3 string
	}{}

	_ = json.Unmarshal([]byte(setting.Data), &dns)

	err := helpers.SetDNSServers(dns.Server1, dns.Server2, dns.Server3)
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	postgresql.AddorUpdateReplication("SYSTEM_DNS", err == nil, errStr)

	log.Println("DNS update check completed")
}

//Extension Extension Replication Loop
func Extension() {
	log.Println("Checking Extensions updates from Liman")
	extensions := postgresql.GetExtensions()

	if len(extensions) < 1 || extensions[0].ID == "" {
		return
	}

	for _, extension := range extensions {
		extensionName := strings.ToLower(extension.Name)
		replicationName := "EXTENSION-" + extension.ID
		replication := postgresql.GetReplication(replicationName)
		if helpers.IsNewer(replication.UpdatedAt, extension.FileUpdateAt) {
			log.Printf("Extension %v is already up to date.", extensionName)
			continue
		}

		log.Println("Download extension " + extension.ID)
		zipPath, err := DownloadExtension(extension.ID)
		if err != nil {
			postgresql.AddorUpdateReplication(replicationName, true, err.Error())
			continue
		}

		targetPath := "/tmp/" + helpers.RandomStr(16) + "/"
		err = helpers.Unzip(zipPath, targetPath)
		if err != nil {
			postgresql.AddorUpdateReplication(replicationName, true, err.Error())
			continue
		}

		//Github direct download fix.
		files, _ := ioutil.ReadDir(targetPath)
		if len(files) == 1 {
			targetPath = targetPath + files[0].Name()
		}

		log.Printf("Extension %v is downloaded and extracted\n", extension.ID)

		extensionFolder := helpers.ExtensionsPath + extensionName

		if _, err := os.Stat(extensionFolder); os.IsNotExist(err) {
			err = os.MkdirAll(extensionFolder, 0700)
			if err != nil {
				postgresql.AddorUpdateReplication(replicationName, true, err.Error())
				continue
			}
		}

		err = helpers.CopyFolder(targetPath, extensionFolder)

		if err != nil {
			postgresql.AddorUpdateReplication(replicationName, true, err.Error())
			continue
		}

		cleanExtensionID := strings.ReplaceAll(extension.ID, "-", "")
		userExists, err := helpers.ExecuteCommand("id " + cleanExtensionID + " >/dev/null 2>&1 && echo 1 || echo 0")

		if userExists == "0" {
			helpers.AddUser(extension.ID)
		}

		helpers.FixExtensionPermissions(extension.ID, extensionName)

		keyPath := helpers.KeysPath + extension.ID
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			randomKey := []byte(helpers.RandomStr(16))
			err := ioutil.WriteFile(helpers.KeysPath+extension.ID, randomKey, 0644)
			if err != nil {
				postgresql.AddorUpdateReplication(replicationName, true, err.Error())
				continue
			}
		}

		helpers.FixExtensionKeys(extension.ID)

		//Check Dependencies
		databasePath := extensionFolder + "/db.json"
		data, err := ioutil.ReadFile(databasePath)
		if err != nil {
			postgresql.AddorUpdateReplication(replicationName, true, err.Error())
			continue
		}

		var database map[string]interface{}
		json.Unmarshal(data, &database)

		if dependencies, ok := database["dependencies"]; ok {
			depStr := fmt.Sprintf("%v", dependencies)
			//Install them.
			command := "if [ -z '$(find /var/cache/apt/pkgcache.bin -mmin -60)' ]; then sudo apt-get update; fi;DEBIAN_FRONTEND=noninteractive sudo apt-get install -o Dpkg::Use-Pty=0 -o Dpkg::Options::='--force-confdef' -o Dpkg::Options::='--force-confold' " + depStr + " -y --force-yes"
			helpers.ExecuteWithLiveOutput(command)
		}

		postgresql.AddorUpdateReplication(replicationName, true, "")

	}

	log.Println("Extensions update check completed")
}

//ExtensionDB ExtensionDB Replication Loop
func ExtensionDB() {
	log.Println("Checking ExtensionDB's updates from Liman")
	extensions := postgresql.GetExtensions()

	if len(extensions) < 1 || extensions[0].ID == "" {
		return
	}

	for _, extension := range extensions {
		extensionName := strings.ToLower(extension.Name)
		replicationName := "EXTENSION_DB-" + extension.ID
		replication := postgresql.GetReplication(replicationName)
		if helpers.IsNewer(replication.UpdatedAt, extension.FileUpdateAt) {
			log.Printf("Extension %v is already up to date.", extensionName)
			continue
		}
		setting := postgresql.GetSystemSetting("EXTENSION_DB-" + extension.ID)

		if setting.ID == "" {
			log.Printf("Extension %v db not found!", extensionName)
			continue
		}
		var database map[string]interface{}
		json.Unmarshal([]byte(setting.Data), &database)
		strDB, _ := json.Marshal(database)
		err := ioutil.WriteFile("/tmp/mert", strDB, 0644)

		postgresql.AddorUpdateReplication(replicationName, err == nil, "")

	}

	log.Println("Extensions update check completed")
}

//Certificate Certificate Replication Loop
func Certificate() {
	log.Println("Checking Certificate updates from Liman")
	setting := postgresql.GetSystemSetting("SYSTEM_CERTIFICATES")
	if setting.ID == "" {
		return
	}

	replication := postgresql.GetReplication("SYSTEM_CERTIFICATES")

	if helpers.IsNewer(replication.UpdatedAt, setting.UpdatedAt) {
		log.Println("Certificates already up to date.")
		return
	}

	var certificates = []struct {
		Certificate string `json:"certificate"`
		TargetName  string `json:"targetName"`
	}{}
	_ = json.Unmarshal([]byte(setting.Data), &certificates)

	errStr := ""
	var err error

	certPath, _ := helpers.GetCertificateStrings()

	helpers.ExecuteCommand("rm " + certPath + "liman*.crt")
	for _, certificate := range certificates {
		err = helpers.AddSystemCertificate(certificate.Certificate, certificate.TargetName)
		if err != nil {
			errStr = err.Error()
			break
		}
	}

	postgresql.AddorUpdateReplication("SYSTEM_CERTIFICATES", err == nil, errStr)

	log.Println("Certificate update check completed")
}

//Loop Main Loop Function.
func Loop() {
	for {
		DNS()
		Extension()
		ExtensionDB()
		Certificate()
		time.Sleep(time.Second * 5)
	}
}
