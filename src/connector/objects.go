package connector

import (
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/masterzen/winrm"
	"github.com/pkg/sftp"
	"github.com/rgzr/sshtun"
	"golang.org/x/crypto/ssh"
)

//Connection Connection Struct
type Connection struct {
	SSH            *ssh.Client
	SFTP           *sftp.Client
	SMB            *smb2.Session
	WinRM          *winrm.Client
	LastConnection time.Time
	WindowsLetter  string
	WindowsPath    string
}

//ActiveConnections Active Connections
var ActiveConnections map[string]Connection

//CloseAllConnections CloseAllConnections
func CloseAllConnections(obj Connection) {
	if obj.SSH != nil {
		obj.SSH.Close()
	}

	if obj.SFTP != nil {
		obj.SFTP.Close()
	}

	if obj.SMB != nil {
		obj.SMB.Logoff()
	}
}

//ActiveTunnels ActiveTunnels
var ActiveTunnels map[string]Tunnel

//Tunnel Tunnel Struct
type Tunnel struct {
	Tunnel         *sshtun.SSHTun
	Port           int
	LastConnection time.Time
}
