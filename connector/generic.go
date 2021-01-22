package connector

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/limanmys/go/postgresql"

	"golang.org/x/crypto/ssh"
	"golang.org/x/text/encoding/unicode"
)

//GetConnection GetConnection
func GetConnection(userID string, serverID string, IPAddress string) (*Connection, error) {
	var val Connection
	if val2, ok := ActiveConnections[userID+serverID+IPAddress]; ok {
		val = val2
	} else {
		res := val.CreateShell(userID, serverID, IPAddress)
		if res == false {
			return &val, errors.New("cannot connect to server")
		}
	}

	val.LastConnection = time.Now()
	ActiveConnections[userID+serverID+IPAddress] = val
	return &val, nil
}

//GetConnectionRaw GetConnectionRaw
func GetConnectionRaw(userID string, connectionType string, username string, password string, IPAddress string, keyPort string) (*Connection, error) {
	var val Connection
	if val2, ok := ActiveConnections[userID+username+IPAddress]; ok {
		val = val2
	} else {
		res := val.CreateShellRaw(connectionType, username, password, IPAddress, keyPort)
		if res == false {
			return &val, errors.New("cannot connect to server")
		}
	}

	val.LastConnection = time.Now()
	ActiveConnections[userID+username+IPAddress] = val
	return &val, nil
}

//CreateShell CreateShell
func (val *Connection) CreateShell(userID string, serverID string, IPAddress string) bool {
	username, password, keyPort, keyObject := postgresql.GetServerKey(userID, serverID)
	if keyObject.Type == "ssh" {
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if keyObject.Type == "ssh_certificate" {
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if keyObject.Type == "winrm" {
		connection, err := InitWinRMShell(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.WinRM = connection
		val.WindowsPath = strings.TrimSpace(val.Run("echo $env:TEMP")) + "\\"
		val.WindowsLetter = val.WindowsPath[0:1]
		val.WindowsPath = val.WindowsPath[3:]
	} else {
		return false
	}
	return true
}

//CreateFileConnection CreateFileConnection
func (val *Connection) CreateFileConnection(userID string, serverID string, IPAddress string) bool {
	username, password, _, keyObject := postgresql.GetServerKey(userID, serverID)
	if keyObject.Type == "ssh" || keyObject.Type == "ssh_certificate" {
		if val.SFTP != nil {
			return true
		}

		flag := val.CreateShell(userID, serverID, IPAddress)
		if flag == false {
			return false
		}
		val.SFTP = OpenSFTPConnection(val.SSH)
	} else if keyObject.Type == "winrm" {
		if val.SMB != nil {
			return true
		}
		temp, err := OpenSMBConnection(IPAddress, username, password)
		if err != nil {
			return false
		}
		val.SMB = temp
	} else {
		return false
	}
	return true
}

//CreateShellRaw CreateShellRaw
func (val *Connection) CreateShellRaw(connectionType string, username string, password string, IPAddress string, keyPort string) bool {
	if connectionType == "ssh" {
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if connectionType == "ssh_certificate" {
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if connectionType == "winrm" {
		connection, err := InitWinRMShell(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.WinRM = connection
	} else {
		return false
	}
	return true
}

//Run Run through ssh
func (val Connection) Run(command string) string {
	if val.SSH != nil {
		sess, _ := val.SSH.NewSession()
		defer sess.Close()
		output, _ := sess.Output(command)
		return strings.TrimSpace(string(output))
	} else if val.WinRM != nil {
		command = "$ProgressPreference = 'SilentlyContinue';" + command
		encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		encoded, _ := encoder.String(command)
		command = base64.StdEncoding.EncodeToString([]byte(encoded))
		stdout, stderr, _, err := val.WinRM.RunWithString("powershell.exe -encodedCommand "+command, "")
		if err != nil {
			return err.Error()
		}
		return strings.TrimSpace(stdout) + strings.TrimSpace(stderr)
	}
	return "Cannot run command!"
}

//Stdin Stdin
func (val *Connection) Stdin() (io.WriteCloser, error) {
	if val.SSH != nil {
		if val.SSHSesion == nil {
			sess, err := val.SSH.NewSession()
			if err != nil {
				fmt.Println(err.Error())
			}
			val.SSHSesion = sess
		}
		pipe, err := val.SSHSesion.StdinPipe()
		val.In = pipe
		return val.In, err
	}
	return nil, errors.New("unsupported connection type")
}

//Magic Magic
func (val *Connection) Magic() error {
	if val.SSH != nil {
		if val.SSHSesion == nil {
			val.SSHSesion, _ = val.SSH.NewSession()
		}
		pipe, err := val.SSHSesion.StdoutPipe()
		val.Out = pipe

		if err != nil {
			return err
		}

		pipe2, err := val.SSHSesion.StdinPipe()
		val.In = pipe2

		if err != nil {
			return err
		}

		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		val.SSHSesion.RequestPty("xterm", 400, 256, modes)
		val.SSHSesion.Start("/bin/bash")

		return nil
	}
	return errors.New("unsupported connection type")
}

//Put Put
func (val Connection) Put(localPath string, remotePath string) bool {
	if val.SFTP != nil {
		return PutFileSFTP(val.SFTP, localPath, remotePath)
	} else if val.SMB != nil {
		return PutFileSMB(val.SMB, localPath, remotePath, val.WindowsLetter)
	}
	return false
}

//Get Get
func (val Connection) Get(localPath string, remotePath string) bool {
	if val.SFTP != nil {
		return GetFileSFTP(val.SFTP, localPath, remotePath)
	} else if val.SMB != nil {
		return GetFileSMB(val.SMB, localPath, remotePath, val.WindowsLetter)
	}
	return false
}

//VerifyAuth VerifyAuth
func VerifyAuth(username string, password string, ipAddress string, port string, keyType string) bool {
	if keyType == "ssh" {
		return VerifySSH(username, password, ipAddress, port)
	} else if keyType == "ssh_certificate" {
		return VerifySSHCertificate(username, password, ipAddress, port)
	} else if keyType == "winrm" {
		return VerifyWinRM(username, password, ipAddress, port)
	}
	return true
}
