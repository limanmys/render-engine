package connector

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/limanmys/go/postgresql"

	"golang.org/x/text/encoding/unicode"

	"github.com/acarl005/stripansi"

	"golang.org/x/crypto/ssh"
)

//GetConnection GetConnection
func GetConnection(userID string, serverID string, IPAddress string) (*Connection, error) {
	var val Connection
	if val2, ok := ActiveConnections[userID+serverID]; ok {
		val = val2
	} else {
		res := val.CreateShell(userID, serverID, IPAddress)
		if res == false {
			return &val, errors.New("cannot connect to server")
		}
	}

	val.LastConnection = time.Now()
	ActiveConnections[userID+serverID] = val
	return &val, nil
}

//CreateShell CreateShell
func (val *Connection) CreateShell(userID string, serverID string, IPAddress string) bool {
	username, password, keyPort, keyObject := postgresql.GetServerKey(userID, serverID)
	val.password = password
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
	} else if keyObject.Type == "winrm_insecure" {
		connection, err := InitInsecureWinRMShell(username, password, IPAddress, keyPort)
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
	val.password = password
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
	} else if keyObject.Type == "winrm_insecure" {
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
	val.password = password
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
	} else if connectionType == "winrm_insecure" {
		connection, err := InitInsecureWinRMShell(username, password, IPAddress, keyPort)
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
		closed := false
		defer sess.Close()
		defer func(closed *bool) {
			*closed = true
		}(&closed)
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,     // disable echoing
			ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
			ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		}
		err := sess.RequestPty("dumb", 1000, 1000, modes)
		if err != nil {
			return err.Error()
		}
		stdoutB := new(bytes.Buffer)
		sess.Stdout = stdoutB
		in, err := sess.StdinPipe()
		if err != nil {
			return err.Error()
		}
		go func(in io.Writer, output *bytes.Buffer, closed *bool) {
			for {
				if *closed {
					break
				}
				if output.Len() > 0 {
					if output.String() == "liman-pass-sudo" {
							_, _ = in.Write([]byte(val.password + "\n"))
						break
					} else {
						break
					}
				}
			}
		}(in, stdoutB, &closed)
		sess.Run("(" + command + ") 2> /dev/null")
		return stripansi.Strip(strings.TrimSpace(strings.Replace(stdoutB.String(), "liman-pass-sudo", "", 1)))
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
	} else if keyType == "winrm_insecure" {
		return VerifyInsecureWinRM(username, password, ipAddress, port)
	}
	return true
}
