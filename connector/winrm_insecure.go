package connector

import (
	"strconv"
	"strings"

	"github.com/masterzen/winrm"
)

//InitInsecureWinRMShell InitInsecureWinRMShell
func InitInsecureWinRMShell(username string, password string, ipAddress string, port string) (*winrm.Client, error) {
	intPort, _ := strconv.Atoi(port)
	endpoint := winrm.NewEndpoint(ipAddress, intPort, false, true, nil, nil, nil, 0)

	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	client, err := winrm.NewClientWithParameters(endpoint, username, password, params)
	if err != nil {
		return nil, err
	}
	return client, nil

}

//VerifyInsecureWinRM VerifyInsecureWinRM
func VerifyInsecureWinRM(username string, password string, ipAddress string, port string) bool {
	intPort, _ := strconv.Atoi(port)
	endpoint := winrm.NewEndpoint(ipAddress, intPort, false, true, nil, nil, nil, 0)

	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	client, err := winrm.NewClientWithParameters(endpoint, username, password, params)
	if err != nil {
		return false
	}

	stdout, _, _, _ := client.RunWithString("hostname", "")
	if strings.TrimSpace(stdout) == "" {
		return false
	}
	return true
}
