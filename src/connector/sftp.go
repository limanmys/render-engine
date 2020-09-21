package connector

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//OpenSFTPConnection Open SFTP connection through active shell
func OpenSFTPConnection(conn *ssh.Client) *sftp.Client {
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal("SFTP > " + err.Error())
	}
	return client
}

//PutFileSFTP Send file
func PutFileSFTP(conn *sftp.Client, localPath string, remotePath string) bool {
	w := conn.Walk(filepath.Dir(remotePath))
	for w.Step() {
		if w.Err() != nil {
			continue
		}
	}

	f, err := conn.Create(filepath.Base(remotePath))
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

//GetFileSFTP Send file
func GetFileSFTP(conn *sftp.Client, localPath string, remotePath string) bool {
	w := conn.Walk(filepath.Dir(remotePath))
	for w.Step() {
		if w.Err() != nil {
			continue
		}
	}

	f, err := conn.Open(remotePath)
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
		fmt.Println(err.Error())
		return false
	}
	return true
}
