package connector

import (
	"github.com/masterzen/winrm"
)

//InitWinRMShell InitWinRMShell
func InitWinRMShell(username string, password string, ipAddress string, port string) (*winrm.Client, error) {
	endpoint := winrm.NewEndpoint(ipAddress, 5986, true, true, nil, nil, nil, 0)

	params := winrm.DefaultParameters
	params.TransportDecorator = func() winrm.Transporter { return &winrm.ClientNTLM{} }

	client, err := winrm.NewClientWithParameters(endpoint, "mcelen", "Passw0rd", params)
	if err != nil {
		return nil, err
	}
	return client, nil

}
