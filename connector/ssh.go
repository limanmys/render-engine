package connector

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

//InitShellWithPassword Initialize shell
func InitShellWithPassword(username string, password string, hostname string, port string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	ipAddress := hostname
	for i := 0; i < 10; i++ {
		addr, err := net.LookupIP(hostname)
		if err == nil {
			ipAddress = addr[0].String()
			break
		}
	}

	conn, err := ssh.Dial("tcp", ipAddress+":"+port, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//InitShellWithCertificate Initialize shell with certificate
func InitShellWithCertificate(username string, certificate string, hostname string, port string) (*ssh.Client, error) {
	key, err := ssh.ParsePrivateKey([]byte(certificate))
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//VerifySSH VerifySSH
func VerifySSH(username string, password string, ipAddress string, port string) bool {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}
	conn, err := ssh.Dial("tcp", ipAddress+":"+port, config)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

//VerifySSHCertificate VerifySSHCertificate
func VerifySSHCertificate(username string, certificate string, ipAddress string, port string) bool {
	key, err := ssh.ParsePrivateKey([]byte(certificate))
	if err != nil {
		return false
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	}

	conn, err := ssh.Dial("tcp", ipAddress+":"+port, config)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}
