package constants

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//Connection Connection Struct
type Connection struct {
	SSH  *ssh.Client
	SFTP *sftp.Client
}

//ActiveConnections Active Connections
var ActiveConnections map[string]Connection

func closeAllConnections(obj Connection) {
	if obj.SSH != nil {
		obj.SSH.Close()
	}

	if obj.SFTP != nil {
		obj.SFTP.Close()
	}
}
