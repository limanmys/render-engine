package connector

import (
	"fmt"
	"renderer/src/constants"
	"renderer/src/sqlite"
)

//CreateConnection CreateConnection
func CreateConnection(userID string, serverID string, IPAddress string, val *constants.Connection) bool {
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

//RawCreateConnection RawCreateConnection
func RawCreateConnection(val *constants.Connection, connectionType string, username string, password string, IPAddress string, keyPort string) bool {
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

//RunCommand Run Command through ssh
func RunCommand(val *constants.Connection, command string) string {
	if val.SSH != nil {
		sess, err := val.SSH.NewSession()
		defer sess.Close()
		output, err := sess.Output(command + " 2>&1")
		if err != nil {
			return ""
		}
		return string(output)
	} else if val.WinRM != nil {
		stdout, stderr, _, err := val.WinRM.RunWithString(command, "")
		if err != nil {
			fmt.Println(err.Error())
		}

		return stdout + stderr
	}
	return ""
}
