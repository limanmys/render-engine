package replications

import (
	"encoding/json"
	"log"
	"time"

	"github.com/limanmys/go/helpers"
	"github.com/limanmys/go/postgresql"
)

func Dns() {
	log.Println("Checking DNS updates from Liman")
	setting := postgresql.GetSystemSetting("SYSTEM_DNS")
	if setting.ID == "" {
		return
	}

	/*replication := postgresql.GetReplication("SYSTEM_DNS")

	if helpers.IsNewer(replication.UpdatedAt, setting.UpdatedAt) {
		return
	}*/

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

func Extension() {

}

func Certificate() {
	log.Println("Checking Certificate updates from Liman")
	setting := postgresql.GetSystemSetting("SYSTEM_CERTIFICATES")
	if setting.ID == "" {
		return
	}

	/*replication := postgresql.GetReplication("SYSTEM_CERTIFICATES")

	if helpers.IsNewer(replication.UpdatedAt, setting.UpdatedAt) {
		return
	}*/

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

func Loop() {
	for {
		Dns()
		Extension()
		Certificate()
		time.Sleep(time.Second * 30)
	}
}
