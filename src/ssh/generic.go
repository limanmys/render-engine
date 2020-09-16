package ssh

import (
	"fmt"
	"renderer/src/sqlite"

	"golang.org/x/crypto/ssh"
)

//CreateConnection CreateConnection
func CreateConnection(userID string, serverID string, IPAddress string) (*ssh.Client, bool) {
	var conn *ssh.Client
	username, password, keyPort, keyObject := sqlite.GetServerKey(userID, serverID)
	if keyObject.Type == "ssh" {
		fmt.Println("Creating a new SSH Connection for " + IPAddress)
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return nil, false
		}
		conn = connection
	} else if keyObject.Type == "ssh_certificate" {
		fmt.Println("Creating a new SSH Certificate Connection for " + IPAddress)
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return nil, false
		}
		conn = connection
	}
	return conn, true
}

//RawCreateConnection RawCreateConnection
func RawCreateConnection(connectionType string, username string, password string, IPAddress string, keyPort string) (*ssh.Client, bool) {
	var conn *ssh.Client
	if connectionType == "ssh" {
		fmt.Println("Creating a new SSH Connection for " + IPAddress)
		connection, err := InitShellWithPassword(username, password, IPAddress, keyPort)
		if err != nil {
			return nil, false
		}
		conn = connection
	} else if connectionType == "ssh_certificate" {
		fmt.Println("Creating a new SSH Certificate Connection for " + IPAddress)
		connection, err := InitShellWithCertificate(username, password, IPAddress, keyPort)
		if err != nil {
			return nil, false
		}
		conn = connection
	} else {
		return nil, false
	}
	return conn, true
}
