package ssh

import (
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
	}
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)
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
	}
	conn, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

//RunCommand Run Command through ssh
func RunCommand(conn *ssh.Client, command string) string {
	sess, err := conn.NewSession()
	defer sess.Close()
	output, err := sess.Output(command + " 2>&1")
	if err != nil {
		panic(err)
	}
	return string(output)
}
