package connector

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
)

//OpenSMBConnection OpenSMBConnection
func OpenSMBConnection(ipAddress string, username string, password string) (*smb2.Session, error) {
	dialer := net.Dialer{Timeout: time.Second * 5}
	conn, err := dialer.Dial("tcp", ipAddress+":445")
	if err != nil {
		return nil, err
	}

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return nil, err
	}
	return s, nil
}

//PutFileSMB PutFileSMB
func PutFileSMB(session *smb2.Session, localPath string, remotePath string, remoteDisk string) bool {
	fs, err := session.Mount(remoteDisk + "$")
	if err != nil {
		return false
	}
	defer fs.Umount()

	f, err := fs.Create(remotePath)
	if err != nil {
		return false
	}
	defer f.Close()
	srcFile, err := os.Open(localPath)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	_, err = io.Copy(f, srcFile)
	if err != nil {
		return false
	}
	return true
}

//GetFileSMB GetFileSMB
func GetFileSMB(session *smb2.Session, localPath string, remotePath string, remoteDisk string) bool {
	fs, err := session.Mount(remoteDisk + "$")
	if err != nil {
		return false
	}
	defer fs.Umount()

	f, err := fs.Open(remotePath)
	if err != nil {
		return false
	}

	defer f.Close()

	_, err = os.Stat(localPath)
	if os.IsNotExist(err) {
		os.Create(localPath)
	}

	srcFile, err := os.OpenFile(localPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return false
	}
	defer srcFile.Close()

	_, err = io.Copy(srcFile, f)
	if err != nil {
		return false
	}
	return true
}
