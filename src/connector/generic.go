package connector

import (
	"fmt"
	"renderer/src/sqlite"
)

//CreateShell CreateShell
func (val *Connection) CreateShell(userID string, serverID string, IPAddress string) bool {
	username, password, keyPort, keyObject := sqlite.GetServerKey(userID, serverID)
	if keyObject.Type == "ssh" {
		fmt.Println("Creating a new SSH Connection for " + IPAddress)
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if keyObject.Type == "ssh_certificate" {
		fmt.Println("Creating a new SSH Certificate Connection for " + IPAddress)
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if keyObject.Type == "winrm" {
		fmt.Println("Creating a new WinRM Connection for " + IPAddress)
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

//CreateFileConnection CreateFileConnection
func (val *Connection) CreateFileConnection(userID string, serverID string, IPAddress string) bool {
	username, password, _, keyObject := sqlite.GetServerKey(userID, serverID)
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
		fmt.Println("Creating a new SSH Connection for " + IPAddress)
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if connectionType == "ssh_certificate" {
		fmt.Println("Creating a new SSH Certificate Connection for " + IPAddress)
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return false
		}
		val.SSH = connection
	} else if connectionType == "winrm" {
		fmt.Println("Creating a new WinRM Connection for " + IPAddress)
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
		sess, err := val.SSH.NewSession()
		defer sess.Close()
		output, err := sess.Output(command + " 2>&1")
		if err != nil {
			return string(output) + err.Error()
		}
		return string(output)
	} else if val.WinRM != nil {
		stdout, stderr, _, err := val.WinRM.RunWithString(command, "")
		if err != nil {
			return err.Error()
		}

		return stdout + stderr
	}
	return "Cannot run command!"
}

//Put Put through ssh
func (val Connection) Put(localPath string, remotePath string) bool {
	if val.SFTP != nil {
		return PutFileSFTP(val.SFTP, localPath, remotePath)
	} else if val.SMB != nil {
		return PutFileSMB(val.SMB, localPath, remotePath)
	}
	return false
}

//Get Get through ssh
func (val Connection) Get(localPath string, remotePath string) bool {
	if val.SFTP != nil {
		return GetFileSFTP(val.SFTP, localPath, remotePath)
	} else if val.SMB != nil {
		return GetFileSMB(val.SMB, localPath, remotePath)
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
